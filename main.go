package main

import (
    "fmt"
    "net/http"
    "log"
    "flag"
    "strings"
    "encoding/json"
    "html/template"
    "github.com/julienschmidt/httprouter"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/conf"
    "github.com/synw/microb/events"
    "github.com/synw/microb/datatypes"
    "github.com/synw/microb/db"
)


var Verbosity = flag.Int("v", 1, "Verbosity level")

var Config = conf.GetConf()

var view = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))

func renderTemplate(response http.ResponseWriter, page *datatypes.Page) {
    err := view.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}

func getPage(url string) *datatypes.Page {
	hasSlash := strings.HasSuffix(url, "/")
	index_url := url
	found := false
	var data map[string]interface{}
	page := &datatypes.Page{Url:"404", Title:"Page not found", Content:"Page not found"}
	if (hasSlash == false) {
		index_url = url+"/"
	}
	// remove url mask
	index_url = strings.Replace(index_url,"/x","",-1)
	// hit db
	data, found = db.GetFromUrl(index_url)
	if (found == false) {
		return page
	}
	fields := data["fields"].(map[string]interface{})
	content := fields["content"].(string)
	title := fields["title"].(string)
	page = &datatypes.Page{Url: url, Title: title, Content: content}
	return page
}

func serveRequest(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")
	msg := "Page "+url+" requested"
	event := *events.NewEvent("runtime_info", "http_server", msg)
    events.Handle(&event)
    if strings.HasPrefix(url, "/x/") {
    	fmt.Println("API call", url)
    	page := getPage(url)
	    if (page.Url == "404") {
	    	msg = "404 page not found in database: "+url
	    	event = *events.NewEvent("error", "http_server", msg)
	    	events.Handle(&event)
	    }
		json_bytes, _ := json.Marshal(page)
		fmt.Fprintf(response, "%s\n", json_bytes)
    } else {
    	page := &datatypes.Page{Url: url, Title: "Page not found", Content: "Page not found"}
    	renderTemplate(response, page)
    }
}

func main() {
    router := httprouter.New()
    router.GET("/*url", serveRequest)
    server := conf.GetServer()
    database := conf.GetMainDatabase()
    loc := server.Host+":"+server.Port
    if (*Verbosity > 0) {
		msg := "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		event := events.NewEvent("runtime_info", "http_server", msg)
		events.Handle(event)
	}
	
    log.Fatal(http.ListenAndServe(loc, router))
}

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
var static_url = Config["staticfiles_host"].(string)

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
	if url == "" {
		url = "/"
	}
	msg := url+" page request"
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	status := http.StatusOK
	event := &datatypes.Event{"request", "http_server", msg, d}
    events.Handle(event, *Verbosity)
    page := &datatypes.Page{Url: url, Title: "", Content: ""}
    response = httpResponseWriter{response, &status}
    renderTemplate(response, page)
}

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func serveApi(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")
	msg := url+" api request"
	d := make(map[string]interface{})
	d["status_code"] = http.StatusOK
	event := &datatypes.Event{"request", "http_server", msg, d}
    events.Handle(event, *Verbosity)
    status := http.StatusOK
	page := getPage(url)
    if (page.Url == "404") {
    	d["status_code"] = http.StatusNotFound
    	msg = url+" document not found in the database"
    	event := &datatypes.Event{"request_error", "http_server", msg, d}
    	events.Handle(event, *Verbosity)
    }
	json_bytes, _ := json.Marshal(page)
	response = httpResponseWriter{response, &status}
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func main() {
    router := httprouter.New()
    router.GET("/", serveRequest)
    router.GET("/p/*url", serveRequest)
    router.GET("/x/*url", serveApi)
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    server := conf.GetServer()
    database := conf.GetMainDatabase()
    loc := server.Host+":"+server.Port
    if (*Verbosity > 0) {
		msg := "Server started on "+loc+" for domain "+skittles.BoldWhite(server.Domain)
		msg = msg+" with "+database.Type
		msg = msg+" ("+database.Host+")"
		event := events.NewEvent("runtime_info", "http_server", msg)
		events.Handle(event, *Verbosity)
	}
	
    log.Fatal(http.ListenAndServe(loc, router))
}

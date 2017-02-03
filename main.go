package main 

import (
    "fmt"
    "log"
    "strings"
    "flag"
    "time"
    "net/http"
    "html/template"
    "encoding/json"
    "microb/conf"
    "github.com/synw/microb/db/rethinkdb"
    "github.com/acmacalister/skittles"
)


type Page struct {
    Url string
    Title string
    Content  string
}

var Config = conf.GetConf()
var edit_mode = flag.Bool("e", false, "Enable edit mode")

func getTime() string {
	t := time.Now()
	return t.Format("15:04:05")
}

func getPage(url string) *Page {
	hasSlash := strings.HasSuffix(url, "/")
	index_url := url
	found := false
	var data map[string]interface{}
	page := Page{Url:"404", Title:"Page not found", Content:"Page not found"}
	if (hasSlash == false) {
		index_url = url+"/"
	}
	// remove url mask
	index_url = strings.Replace(index_url,"/x","",-1)
	// hit db
	if (Config["db_type"] == "rethinkdb") {
		data, found = rethinkdb.GetFromDb(index_url)
		if (found == false) {
			return &page
		}
	} else {
		log.Fatal("No db configured")
	}
	fields := data["fields"].(map[string]interface{})
	content := fields["content"].(string)
	title := fields["title"].(string)
	page = Page{Url: url, Title: title, Content: content}
	return &page
}

var view = template.Must(template.New("view.html").ParseFiles("view.html", "routes.js"))

func renderTemplate(response http.ResponseWriter, page *Page) {
    err := view.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}

func viewHandler(response http.ResponseWriter, request *http.Request) {
    url := request.URL.Path
    fmt.Printf("%s Page %s\n", getTime(), url)
    page := &Page{Url: url, Title: "Page not found", Content: "Page not found"}
    renderTemplate(response, page)
}

func apiHandler(response http.ResponseWriter, request *http.Request) {
    url := request.URL.Path
    page := getPage(url)
    status := ""
    if (page.Url == "404") {
		status = "Error 404"
	}
    fmt.Printf("%s API %s %s\n", getTime(), url, skittles.BoldRed(status))
	json_bytes, _ := json.Marshal(page)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func main() {
	flag.Parse()
	if ( *edit_mode == true ) {
		fmt.Println(skittles.Yellow("Warning"), ": edit mode is enabled")
	}
	http_port := Config["http_port"].(string)
	fmt.Printf("Server started on %s...\n", http_port)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media/"))))
	http.HandleFunc("/x/", apiHandler)
    http.HandleFunc("/", viewHandler)
    log.Fatal(http.ListenAndServe(http_port, nil))
}

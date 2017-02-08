package main 

import (
    "fmt"
    "log"
    "strings"
    "flag"
    "net/http"
    "html/template"
    "encoding/json"
    "sync"
    "io/ioutil"
    "time"
    _ "expvar"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/conf"
    "github.com/synw/microb/utils"
    "github.com/synw/microb/db/rethinkdb"
)


type Page struct {
    Url string
    Title string
    Content  string
}

var Nochangefeed = flag.Bool("nf", false, "Do not use changefeed")
var CommandFlag = flag.String("c", "", "Fire command")

var Config = conf.GetConf()

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
	data, found = rethinkdb.GetFromDb(index_url)
	if (found == false) {
		return &page
	}
	fields := data["fields"].(map[string]interface{})
	content := fields["content"].(string)
	title := fields["title"].(string)
	page = Page{Url: url, Title: title, Content: content}
	return &page
}

var view = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))

func renderTemplate(response http.ResponseWriter, page *Page) {
    err := view.Execute(response, page)
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
    }
}

func viewHandler(response http.ResponseWriter, request *http.Request) {
    url := request.URL.Path
    fmt.Printf("%s Page %s\n", utils.GetTime(), url)
    page := &Page{Url: url, Title: "Page not found", Content: "Page not found"}
    renderTemplate(response, page)
}

func apiHandler(response http.ResponseWriter, request *http.Request) {
    url := request.URL.Path
    page := getPage(url)
    fmt.Printf("%s API %s\n", utils.GetTime(), url)
    if (page.Url == "404") {
    	msg := "404 page not found in database: "+url
    	utils.PrintEvent("error", msg)
    }
	json_bytes, _ := json.Marshal(page)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func reparseStatic() {
	utils.PrintEvent("command", "Reparsing static files")
	view = template.Must(template.New("view.html").ParseFiles("templates/view.html", "templates/routes.js"))
}

func updateRoutes(c chan bool) {
	var routestab []string
	// hit db
	routestab = rethinkdb.GetRoutes()
	var routestr string
	var route string
	for i := range(routestab) {
		route = routestab[i]
		routestr = routestr+fmt.Sprintf("page('%s', function(ctx, next) { loadPage('/x%s') } );", route, route)
	}
	utils.PrintEvent("command", "Rebuilding client side routes")
    str := []byte(routestr)
    err := ioutil.WriteFile("templates/routes.js", str, 0644)
    if err != nil {
        panic(err)
    }
	c <- true
}

func init() {
	// changefeed listeners
	c := make(chan *rethinkdb.DataChanges)
	c2 := make(chan bool)
	comchan := make(chan *rethinkdb.Command)
	go rethinkdb.PageChangesListener(c)
	go rethinkdb.CommandsListener(comchan)
	// channel listeners
	go func() {
    	for {
            changes := <- c
			if (changes.Type == "update") {
				utils.PrintEvent("event", changes.Msg)
				go updateRoutes(c2)
			} else if (changes.Type == "delete") {
				utils.PrintEvent("event", changes.Msg)
				go updateRoutes(c2)
			} else if (changes.Type == "insert") {
				utils.PrintEvent("event", changes.Msg)
				go updateRoutes(c2)
			}
    	}
    }()
    go func() {
    	for {
            com := <- comchan
			if (com.Name != "") {
				msg := "Command "+skittles.BoldWhite(com.Name)+" received"
				utils.PrintEvent("event", msg)
				if (com.Name == "reparse_templates") {
					go reparseStatic()
				} 
			}
    	}
    }()
    go func() {
		for {
			routes_done := <- c2
			if (routes_done == true) {
				//fmt.Println("[OK] Routes updated")
				go reparseStatic()
			}
		}
	}()
}

func main() {
	flag.Parse()
	// commands
	if (*CommandFlag != "") {
		valid_commands := []string{"update_routes", "reparse_templates"}
		is_valid := false
		for _, com := range(valid_commands) {
			if (com == *CommandFlag) {
				is_valid = true
				break
			}
		}
		if (is_valid == true) {
			msg := "Sending command "+skittles.BoldWhite(*CommandFlag)+" to the server"
			utils.PrintEvent("event", msg)
			var wg sync.WaitGroup
			wg.Add(1)
			go rethinkdb.SaveCommand(*CommandFlag, &wg)
			wg.Wait()
			
		} else {
			msg := "Unknown command: "+*CommandFlag
			utils.PrintEvent("error", msg)
		}
		return
	}
	if (*Nochangefeed == false) {
		utils.PrintEvent("info", "listening to changefeed")
	}
	// http server
	http_host := Config["http_host"].(string)
	msg := "Server started on "+http_host+" for domain "+Config["domain"].(string)+" with db "+Config["domain"].(string)
	msg = msg+" at "+Config["db_host"].(string)
	utils.PrintEvent("nil", msg)
	server := &http.Server{
		Addr: http_host,
	    ReadTimeout: 5 * time.Second,
	    WriteTimeout: 10 * time.Second,
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/x/", apiHandler)
    http.HandleFunc("/", viewHandler)
    log.Fatal(server.ListenAndServe())
}

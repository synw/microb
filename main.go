package main 

import (
    "fmt"
    "log"
    "strings"
    "strconv"
    "flag"
    "sync"
    "net/http"
    "html/template"
    "encoding/json"
    "time"
    _ "expvar"
    "github.com/acmacalister/skittles"
    "github.com/synw/microb/conf"
    "github.com/synw/microb/utils"
    "github.com/synw/microb/db"
    "github.com/synw/microb/db/datatypes"
    "github.com/synw/microb/middleware"
    "github.com/synw/microb/commands"
)


type Page struct {
    Url string
    Title string
    Content  string
}

var CommandFlag = flag.String("c", "noflag", "Fire command")
var Metrics = flag.Bool("m", false, "Display metrics")
var Verbosity = flag.Int("v", 1, "Verbosity level")
var Reason = flag.String("r", "nil", "Reason for sending a command (to use with the -c flag)")

var Config = conf.GetConf()

var C_display = make(chan string)

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
	data, found = db.GetFromDb(index_url)
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
    purl := request.URL.Path
    //fmt.Printf("%s Page %s\n", utils.GetTime(), url)
    page := &Page{Url: purl, Title: "Page not found", Content: "Page not found"}
    go middleware.ProcessHit(request, Config["hits_log"].(bool), *Verbosity, C_display)
    renderTemplate(response, page)
}

func apiHandler(response http.ResponseWriter, request *http.Request) {
    purl := request.URL.Path
    page := getPage(purl)
    //fmt.Printf("%s API %s\n", utils.GetTime(), url)
    if (page.Url == "404") {
    	msg := "404 page not found in database: "+purl
    	utils.PrintEvent("error", msg)
    }
	json_bytes, _ := json.Marshal(page)
    go middleware.ProcessHit(request, Config["hits_log"].(bool), *Verbosity, C_display)
	fmt.Fprintf(response, "%s\n", json_bytes)
}

func init() {
	flag.Parse()
	c_commands_results := make(chan bool)
	// listen to commands results
    go func() {
		for {
			result := <- c_commands_results
			if (result == true) {
				utils.PrintEvent("nil", "Command successfull")
			} else {
				utils.PrintEvent("error", "Error executing command")
			}
		}
	}()
	// manual commands
	if (*CommandFlag != "noflag") {
		command := &datatypes.Command{*CommandFlag, "terminal", "nil"}
		if *Reason != "nil" {
			command.Reason = *Reason
		}
		var wg sync.WaitGroup
		wg.Add(1)
		commands.RunCommandAndExit(command, &wg, c_commands_results)
		wg.Wait()
	} else {
		if (*Verbosity > 0) {
			if (Config["hits_monitor"] == true) {
				utils.PrintEvent("info", "Monitoring hits")
			}
			if (Config["hits_log"] == true) {
				utils.PrintEvent("info", "Logging hits")
			}
			utils.PrintEvent("info", "Listening to changefeeds")
		}
		// changefeed listeners
		c_pages_changes := make(chan *datatypes.DataChanges)
		comchan := make(chan *datatypes.Command)
		go db.PageChangesListener(c_pages_changes)
		go db.CommandsListener(comchan)
		// listen for change in pages table and trigger handlers
		go func() {
	    	for {
	            changes := <- c_pages_changes
				if (changes.Type == "update") {
					utils.PrintEvent("event", changes.Msg)
					command := &datatypes.Command{"update_routes", "listener", "Update event in the database"}
					go commands.RunCommand(command, c_commands_results, true)
				} /*else if (changes.Type == "delete") {
					utils.PrintEvent("event", changes.Msg)
					go updateRoutes()
				} else if (changes.Type == "insert") {
					utils.PrintEvent("event", changes.Msg)
					go updateRoutes()
				}*/
	    	}
	    }()
	    // listen for incoming commands
	    go func() {
	    	for {
	            com := <- comchan
				if (com.Name != "") {
					msg := "Command "+skittles.BoldWhite(com.Name)+" received"
					utils.PrintEvent("event", msg)
					if (com.Name == "reparse_templates") {
						//go reparseStatic()
					} else if (com.Name == "update_routes") {
						//go updateRoutes(c2)
					} 
				}
	    	}
	    }()
		// hits writer
		c_hits := make(chan int)
		if (Config["hits_monitor"].(bool) == true || Config["hits_log"].(bool) == true) {
			go middleware.WatchHits(1, Config["hits_log"].(bool), Config["hits_monitor"].(bool), c_hits)
		}
		// hits monitor
		if (Config["hits_monitor"].(bool) == true) {
			go func() {
				for {
					num_hits := <- c_hits
					if (num_hits > 0) {
						if (*Metrics == true) {
							msg := "Hits per second: "+strconv.Itoa(num_hits)
							go utils.PrintEvent("metric", msg)
						}
					}
				}
			}()
		}
		// hits display
		if (*Verbosity > 0) {
			go func() {
				for {
					msg := <- C_display
					fmt.Println(msg)
				}
			}()
		}
	}
}

func main() {
	// http server
	http_host := Config["http_host"].(string)
	var msg string
	if (*Verbosity > 0) {
		msg = "Server started on "+http_host+" for domain "+skittles.BoldWhite(Config["domain"].(string))
		msg = msg+" with "+Config["db_type"].(string)+" db "+Config["domain"].(string)
		msg = msg+" ("+Config["db_host"].(string)+")"
	}
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

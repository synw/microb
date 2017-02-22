package rethinkdb

import (
	"fmt"
	"time"
	//"errors"
	 r "gopkg.in/dancannon/gorethink.v3"
	 "github.com/synw/microb/libmicrob/datatypes"
	 "github.com/synw/microb/libmicrob/events"
	 "github.com/synw/microb/libmicrob/state"
)

var conn *r.Session

func InitDb() {
	if state.Server.PagesDb != nil {
		if state.Server.PagesDb.Type == "rethinkdb" {
			conn = connect(state.Server.PagesDb)
		}
	}
}

func connect(database *datatypes.Database) (*r.Session) {
	host := database.Host
	port := database.Port
	user := database.User
	pwd := database.Password
	addr := host+":"+port
	db_name := state.Server.Domain
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		Database: db_name,
		Username: user,
		Password: pwd,
		InitialCap: 10,
        MaxOpen:    10,
	})
    if err != nil {
        events.Error("db.rethinkdb.connectToDb()", err)
    }
    if state.Debug {
    	fmt.Println("Connecting to database", db_name, " at ", addr)
    }
    return session
}
/*
func ReportIssues() []*datatypes.DatabaseIssue {
	session := conn
	res, err := r.Db("rethinkdb").Table("current_issues").Run(session)
	defer res.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	var issues []datatypes.DatabaseIssue
	for res.Next(&issues) {
		if res.Type != "" {
			var issue := *datatypes.DatabaseIssue
		  	issue.Type = res.Type
		  	issue.Description = res.Description
		  	issues = append(issues, issue)
		}
	return issues
}
*/
func ReportStatus()(map[string]interface{}, error) {
	session := conn
	res, err := r.DB("rethinkdb").Table("server_status").Run(session)
	defer res.Close()
	status := make(map[string]interface{})
	if err != nil {
		return status, err
	}
	var row map[string]interface{}
	err = res.One(&row)
	if err != nil && err != r.ErrEmptyResult {
		events.Error("db.rethinkdb.ReportStatus", err)
	}
	process := row["process"].(map[string]interface{})
	//network := row["network"].(map[string]interface{})
	version := process["version"].(string)
	cache_size_mb := process["cache_size_mb"].(float64)
	time_started := process["time_started"].(time.Time)
	final_status := make(map[string]interface{})
	final_status["version"] = version
	final_status["cache_size_mb"] = cache_size_mb
	final_status["time_started"] = time_started
	return final_status, nil
}

func GetRoutes() []string {
	session := conn
	res, err := r.Table("pages").Pluck("uri").Run(session)
	if err != nil && err != r.ErrEmptyResult {
		events.Error("db.rethinkdb.GetRoutes", err)
	}
	defer res.Close()
	var row map[string]interface{}
	var routes []string
	for res.Next(&row) {
		url := row["uri"].(string)
	    routes = append(routes, url)
	}
	/*
	filters := map[string]interface{}{"type":"routes"}
	res, err := r.DB(db).table("metadata").Filter(filters).Run(session)
	if err != nil {
		events.Error("db.rethinkdb.GetRoutes()", err)
	}
	if err == r.ErrEmptyResult {
		msg := errors.New("Routes not found in database: "+err.Error())
	    events.Error("db.rethinkdb.GetRoutes()", msg)
	} else {
		var row map[string]interface{}
		err = res.One(&row)
		if err != nil && err != r.ErrEmptyResult {
			events.Error("db.rethinkdb.GetRoutes()", err)
		} else {
			
		}
		
	}*/
	return routes
}

func GetFromUrl(url string)  (*datatypes.Page, bool, error)  {
	session := conn
	var page datatypes.Page
	filters := map[string]interface{}{"uri":url}
	res, err := r.Table("pages").Filter(filters).Pluck("fields").Run(session)
	defer res.Close()
	var row map[string]interface{}
	err = res.One(&row)
	if (err != nil && err != r.ErrEmptyResult) {
		events.Error("db.rethinkdb.GetFromUrl", err)
		return &page, false, err
	}
	if err == r.ErrEmptyResult {
		msg := "Empty results for url "+url
		events.New("runtime", "db.rethinkdb.GetFromUrl", msg)
		return &page, false, err
	} else {
		fields := row["fields"].(map[string]interface{})
		title := fields["title"].(string)
		content := fields["content"].(string)
		page.Url = url
		page.Title = title
		page.Content = content
		return &page, true, err
	}
	return &page, false, nil
}

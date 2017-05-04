package rethinkdb

import (
	"errors"
	 r "gopkg.in/dancannon/gorethink.v3"
	 "github.com/synw/terr"
	 "github.com/synw/microb/libmicrob/datatypes"
	 "github.com/synw/microb/libmicrob/state"
)


var conn *r.Session
var isConnected bool = false


func GetByUrl(url string)  (*datatypes.Document, bool, *terr.Trace)  {
	session := conn
	var doc datatypes.Document
	filters := map[string]interface{}{"uri":url}
	res, err := r.DB(state.DocDb.Dbs["documents"]).Table(state.DocDb.Tables["documents"]).Filter(filters).Pluck("fields").Run(session)
	defer res.Close()
	var row map[string]interface{}
	err = res.One(&row)
	if (err != nil && err != r.ErrEmptyResult) {
		tr := terr.New("db.rethinkdb.GetFromUrl", err)
		return &doc, false, tr
	}
	if err == r.ErrEmptyResult {
		msg := errors.New("Empty results for url "+url+" - "+err.Error())
		tr := terr.New("db.rethinkdb.GetFromUrl", msg)
		return &doc, false, tr
	} else {
		fields := row["fields"].(map[string]interface{})
		doc.Url = url
		doc.Fields = fields
		return &doc, true, nil
	}
	return &doc, false, nil
}

func InitDb(database *datatypes.Database) *terr.Trace {
	if database.Type != "rethinkdb" {
		err := errors.New("Not a rethinkdb database")
		tr := terr.New("db.rethinkdb.Initdb", err)
		return tr
	}
	cn, err := connect(database)
	conn = cn
	if err != nil {
		tr := terr.New("db.rethinkdb.InitDb", err)
		return tr
	}
	return nil
}
/*
func SwitchDb(role string, database *datatypes.Database) error {
	var err error
	if isConnected == true {
		if state.Debug == true {
			events.Debug("Closing previous connection to rethinkdb")
		}
		conn.Close()
	}
	if role == "pages" {
		if state.Debug == true {
			events.Debug("Initializing new database "+database.Name+" ("+database.Host+") for role "+role)
		}
		err = InitDb(database)
	} else {
		err = errors.New("Not implemented")
		events.Error("db.rethinkdb.SwitchDb", err)
	}
	return err
}*/

func connect(database *datatypes.Database) (*r.Session, *terr.Trace) {
	user := database.User
	pwd := database.Password
	addr := database.Addr
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		Username: user,
		Password: pwd,
		InitialCap: 10,
        MaxOpen:    10,
	})
    if err != nil {
        tr := terr.New("db.rethinkdb.connectToDb()", err)
        return session, tr
    }
    isConnected = true
    return session, nil
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

func ReportStatus()(map[string]interface{}, *terr.Trace) {
	session := conn
	res, err := r.DB("rethinkdb").Table("server_status").Run(session)
	defer res.Close()
	status := make(map[string]interface{})
	if err != nil {
		tr := terr.New("db.rethinkdb.ReportStatus", err)
		return status, tr
	}
	var row map[string]interface{}
	err = res.One(&row)
	if err != nil && err != r.ErrEmptyResult {
		tr := terr.New("db.rethinkdb.ReportStatus", err)
		return status, tr
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
*/

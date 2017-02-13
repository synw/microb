package rethinkdb

import (
	"fmt"
	"log"
	"sync"
	"time"
	"strings"
	 r "gopkg.in/dancannon/gorethink.v2"
	 "github.com/synw/microb/conf"
	 "github.com/synw/microb/db/datatypes"
)

var Config = conf.GetConf()
var Conn *r.Session
var MainDb *datatypes.Database = conf.GetMainDb()

func init() {
	if MainDb.Type == "rethinkdb" {
		Conn = connectToDb()
	}
}

func connectToDb() (*r.Session) {
	db := conf.GetMainDb()
	host := db.Host
	port := db.Port
	user := db.User
	pwd := db.Password
	addr := host+":"+port
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		//Database: db,
		Username: user,
		Password: pwd,
		InitialCap: 10,
        MaxOpen:    10,
	})
    if err != nil {
        log.Fatalln(err.Error())
    }
    return session
}

func MakeHit(doc string) *datatypes.Hit {
	data := strings.Split(doc, "#!#")
	datenow := time.Now()
	hit := &datatypes.Hit{data[0], data[1], data[2], data[3], data[4], datenow}
	return hit
}

func SaveHits(values []string) {
	session := Conn
	db := Config["domain"].(string)
	for _, doc := range values {
		// unpack the data
		hit := MakeHit(doc)
		// package the data in json
	   	_, err := r.DB(db).Table("hits").Insert(hit, r.InsertOpts{Durability: "soft", ReturnChanges: false}).Run(session)
		if err != nil {
			fmt.Println("Rehinkdb: error writing hits data:", err)
		}
	}
}

func GetRoutes() []string {
	session := Conn
	db := Config["domain"].(string)
	res, err := r.DB(db).Table("pages").Pluck("uri").Run(session)
	defer res.Close()
	var row map[string]interface{}
	var routes []string
	for res.Next(&row) {
		url := row["uri"].(string)
	    routes = append(routes, url)
	}
	if err != nil {
		fmt.Printf("Rethinkdb: error scanning database results: %s\n", err)
	}
	return routes
}

func GetAllData(database *datatypes.Database, table string) *r.Cursor  {
	session := Conn
	res, err := r.DB(database.Name).Table(table).GetAll().Run(session)
	defer session.Close()
	if err == r.ErrEmptyResult {
	    fmt.Printf("Rethinkdb: no data in table: %s\n", err)
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	return res
}

func GetFromDb(url string)  (map[string]interface{}, bool)  {
	session := Conn
	found := false
	db := Config["domain"].(string)
	filters := map[string]interface{}{"uri":url}
	res, err := r.DB(db).Table("pages").Filter(filters).Pluck("fields").Run(session)
	defer res.Close()
	if err == r.ErrEmptyResult {
	    fmt.Printf("Rethinkdb: no results for uri scan: %s\n", err)
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	var rescol map[string]interface{}
	err = res.One(&rescol)
	if err != nil {
		fmt.Printf("Rethinkdb: error scanning database results: %s\n", err)
	}
	var page_served map[string]interface{}
	if err == r.ErrEmptyResult {
	    //fmt.Printf("Rethinkdb: no results: %s\n", err)
	    page_served = make(map[string]interface{})
	} else {
		found = true
		page_served = rescol
	}
	return page_served, found
}

func scanForChanges(res map[string]interface{}) *datatypes.DataChanges {
	new_val := res["new_val"].(map[string]interface{})
	old_val := res["old_val"].(map[string]interface{})
	new_domain := new_val["domain"].(string)
	old_domain := old_val["domain"].(string)
	new_uri := new_val["uri"].(string)
	old_uri := old_val["uri"].(string)
	new_editor := new_val["editor"].(string)
	old_editor := old_val["editor"].(string)
	new_fields := new_val["fields"].(map[string]interface{})
	old_fields := old_val["fields"].(map[string]interface{})
	new_title := new_fields["title"].(string)
	old_title := old_fields["title"].(string)
	new_content := new_fields["content"].(string)
	old_content := old_fields["content"].(string)
	msg := ""
	ctype := "nil"
	changes := make(map[string]interface{})
	has_changed := false
	if (new_domain != old_domain) {
		changes["domain"] = []string{new_domain, old_domain}
		has_changed = true
	}
	if (new_uri != old_uri) {
		changes["uri"] = []string{new_uri, old_uri}
		has_changed = true
	}
	if (new_editor != old_editor) {
		changes["editor"] = []string{new_editor, old_editor}
		has_changed = true
	}
	if (new_content != old_content) {
		changes["content"] = []string{new_content, old_content}
		has_changed = true
	}
	if (new_title != old_title) {
		changes["title"] = []string{new_title, old_title}
		has_changed = true
	}
	if (has_changed == true) {
		ctype = "update"
		var fieldsChanges string
		for k, _ := range(changes) {
			fieldsChanges = fieldsChanges+k+" "
		}
		msg = "Changes detected in the database ("+ctype+"):\n"
		m := "-> "+new_uri+": "+new_editor+" changed fields: "+fieldsChanges
		msg = msg+m
	}
	dataChanges := &datatypes.DataChanges{Msg:msg, Type:ctype, Values:changes}
	return dataChanges
}

func SaveCommand(command *datatypes.Command) {
	session := Conn
	defer session.Close()
	db := Config["database"].(string)
	_, err := r.DB(db).Table("commands").Insert(command, r.InsertOpts{ReturnChanges: false}).RunWrite(session)
	if err != nil { log.Fatalln(err) }
}

func SaveCommandWait(command *datatypes.Command, wg *sync.WaitGroup) {
	session := Conn
	defer wg.Done()
	defer fmt.Println("save command wait")
	defer session.Close()
	db := Config["database"].(string)
	_, err := r.DB(db).Table("commands").Insert(command, r.InsertOpts{ReturnChanges: false}).RunWrite(session)
	if err != nil { log.Fatalln(err) }
}

func CommandsListener(comchan chan *datatypes.Command) {
	session := Conn
	db := Config["domain"].(string)
	// monitor commands
	com, err := r.DB(db).Table("commands").Pluck("Name").Changes().Run(session)
	defer com.Close()
	if err != nil { log.Fatalln(err) }
	var commands map[string]interface{}
	for com.Next(&commands) {
		comt := commands["new_val"].(map[string]interface{})
		comc := comt["Name"].(string)
		coms := &datatypes.Command{Name:comc}
		comchan <- coms
	}
}

func PageChangesListener(c chan *datatypes.DataChanges) {
	session := Conn
	db := Config["domain"].(string)
	// monitor changes in pages table
	res, err := r.DB(db).Table("pages").Changes().Run(session)
	defer res.Close()
	if err != nil { log.Fatalln(err) }
	var changes map[string]interface{}
	var dataChanged *datatypes.DataChanges
	defer res.Close()
	for res.Next(&changes) {
		//fmt.Println(&changes)
		var new_val map[string]interface{}
		var old_val map[string]interface{}
		if (changes["old_val"] == nil) {
			new_val = changes["new_val"].(map[string]interface{})
			old_val = nil
		}
		if (changes["new_val"] == nil) {
			old_val = changes["old_val"].(map[string]interface{})
			new_val = nil
		}
		if (changes["new_val"] != nil && changes["old_val"] != nil) {
			new_val = changes["new_val"].(map[string]interface{})
			old_val = changes["old_val"].(map[string]interface{})
		}
		var msg string
		if (new_val != nil && old_val != nil) {
			dataChanged = scanForChanges(changes)
		} else if (old_val == nil) {
			msg = "Changes detected in the database (insert):\n"
			msg = msg+"-> "+new_val["editor"].(string)+" created a new document: "+new_val["uri"].(string)
			dataChanged = &datatypes.DataChanges{Msg:msg, Type:"insert", Values:new_val}
		} else if (new_val == nil) {
			msg = "Changes detected in the database (delete):\n"
			msg = msg+"-> Document "+old_val["uri"].(string)+" has been deleted"
			dataChanged = &datatypes.DataChanges{Msg:msg, Type:"delete", Values:old_val}
		}
		c <- dataChanged
	}
}

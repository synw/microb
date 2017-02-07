package rethinkdb

import (
	"fmt"
	"log"
	"sync"
	 r "gopkg.in/dancannon/gorethink.v2"
	 "github.com/synw/microb/conf"
	 //"microb/utils"
)

var Config = conf.GetConf()

type DataChanges struct {
	Msg string
	Type string
	Values map[string]interface{}
}

type Command struct {
	Name string
}

func connectToDb() (*r.Session) {
	host := Config["db_host"].(string)
	port := Config["db_port"].(string)
	db := Config["database"].(string)
	user := Config["db_user"].(string)
	pwd := Config["db_password"].(string)
	addr := host+":"+port
	// connect to Rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: addr,
		Database: db,
		Username: user,
		Password: pwd,
	})
    if err != nil {
        log.Fatalln(err.Error())
    }
    return session
}

func GetRoutes() []string {
	session := connectToDb()
	db := Config["database"].(string)
	table := Config["table"].(string)
	res, err := r.DB(db).Table(table).Pluck("uri").Run(session)
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

func GetFromDb (url string)  (map[string]interface{}, bool)  {
	session := connectToDb()
	found := false
	db := Config["database"].(string)
	table := Config["table"].(string)
	filters := map[string]interface{}{"uri":url}
	res, err := r.DB(db).Table(table).Filter(filters).Pluck("fields").Run(session)
	if err == r.ErrEmptyResult {
	    fmt.Printf("Rethinkdb: no results for uri scan: %s\n", err)
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer res.Close()
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

func scanForChanges(res map[string]interface{}) *DataChanges {
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
	dataChanges := &DataChanges{Msg:msg, Type:ctype, Values:changes}
	return dataChanges
}

func SaveCommand(command string, wg *sync.WaitGroup) {
	session := connectToDb()
	defer wg.Done()
	defer session.Close()
	db := Config["database"].(string)
	new_command := &Command{Name:command}
	_, err := r.DB(db).Table("commands").Insert(new_command, r.InsertOpts{ReturnChanges: false}).RunWrite(session)
	if err != nil { log.Fatalln(err) }
}

func CommandsListener(comchan chan *Command) {
	session := connectToDb()
	db := Config["database"].(string)
	// monitor commands
	com, err := r.DB(db).Table("commands").Pluck("Name").Changes().Run(session)
	if err != nil { log.Fatalln(err) }
	var commands map[string]interface{}
	for com.Next(&commands) {
		comt := commands["new_val"].(map[string]interface{})
		comc := comt["Name"].(string)
		coms := &Command{Name:comc}
		comchan <- coms
	}
}

func PageChangesListener(c chan *DataChanges) {
	session := connectToDb()
	db := Config["database"].(string)
	table := Config["table"].(string)
	// monitor changes in pages table
	res, err := r.DB(db).Table(table).Changes().Run(session)
	if err != nil { log.Fatalln(err) }
	var changes map[string]interface{}
	var dataChanged *DataChanges
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
			dataChanged = &DataChanges{Msg:msg, Type:"insert", Values:new_val}
		} else if (new_val == nil) {
			msg = "Changes detected in the database (delete):\n"
			msg = msg+"-> Document "+old_val["uri"].(string)+" has been deleted"
			dataChanged = &DataChanges{Msg:msg, Type:"delete", Values:old_val}
		}
		c <- dataChanged
	}
}

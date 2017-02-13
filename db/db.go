package db

import (
	"fmt"
	"sync"
	r "gopkg.in/dancannon/gorethink.v2"
	postgresql "github.com/synw/microb/db/postgresql"
	"github.com/synw/microb/conf"
	"github.com/synw/microb/db/datatypes"
	"github.com/synw/microb/db/rethinkdb"
	//"github.com/synw/microb/utils"
	
)

var Config = conf.GetConf()
var MainDb *datatypes.Database = conf.GetMainDb()
var Backend = MainDb.Type

func getDatabase(dbname string) *datatypes.Database {
	db_in_conf := Config["databases"].(map[string]interface{})[dbname].(map[string]interface{})
	database := &datatypes.Database{db_in_conf["type"].(string), dbname, db_in_conf["host"].(string), db_in_conf["port"].(string), db_in_conf["user"].(string),db_in_conf["password"].(string)}
	return database
}

func SaveCommand(command *datatypes.Command) {
	if Backend == "rethinkdb" {
		rethinkdb.SaveCommand(command)
	}
	return
}

func SaveCommandWait(command *datatypes.Command, wg *sync.WaitGroup) {
	if Backend == "rethinkdb" {
		rethinkdb.SaveCommandWait(command, wg)
	}
	return
}

func SaveHits(values []string) {
	if Backend == "rethinkdb" {
		rethinkdb.SaveHits(values)
	}
	return
}

func GetFromDb(url string)  (map[string]interface{}, bool) {
	res := make(map[string]interface{})
	ok := false
	if Backend == "rethinkdb" {
		res, ok = rethinkdb.GetFromDb(url)
	}
	return res, ok
}

func GetRoutes() []string {
	routes := []string{}
	if Backend == "rethinkdb" {
		routes = rethinkdb.GetRoutes()
	}
	return routes
}

func ImportPagesFromMainDb(dbname string) {
	// get data
	var res *r.Cursor
	if Backend == "rethinkdb" {
		res = rethinkdb.GetAllData(MainDb, "pages")
	}
	//created := 0
	//updated := 0
	// save data
	to_db := getDatabase(dbname)
	if to_db.Type == "postgresql" {
		pgdb := postgresql.Connect(to_db)
		defer pgdb.Close()
		// check for table
		if pgdb.HasTable("pages") == false {
			pgdb.CreateTable(&datatypes.Database{})
		}
		// iterate over results and save or update rows
		var row map[string]interface{}
		for res.Next(&row) {
			domain := row["domain"].(string)
			uri := row["uri"].(string)
			title := row["title"].(string)
			content := row["content"].(string)
			page := &datatypes.GormPage{domain, uri, title, content}
			// write to pg
			pgres := pgdb.Where(&datatypes.GormPage{Url: uri}).Assign(page).FirstOrCreate(&page)
			fmt.Println("res", pgres)
			/*
			if pgdb.NewRecord(page) == true {
				pgdb.Create(page)
				created = created+1
			} else {
				pgpage := pgdb.First(page)
				pgpage.Domain = domain
				pgpage.Uri = uri
				pgpage.Title = title
				pgpage.Content = content
				pgpage.Save()
				updated = updated+1
			}*/
		}
	}
	/*ln := created+" pages created"
	msg := fmt.Println(ln)
	ln = ln+updated+" pages updated"
	msg = msg+fmt.Println(ln)
	utils.PrintEvent("event", msg)*/
	return
}

func CommandsListener(comchan chan *datatypes.Command) {
	if conf.ListenToChangefeeds() == true {
		rethinkdb.CommandsListener(comchan)
	}
}

func PageChangesListener(c chan *datatypes.DataChanges) {
	if conf.ListenToChangefeeds() == true {
		rethinkdb.PageChangesListener(c)
	}
}

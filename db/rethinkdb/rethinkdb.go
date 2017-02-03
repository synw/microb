package rethinkdb

import (
	"fmt"
	"log"
	 r "gopkg.in/dancannon/gorethink.v2"
	 "github.com/synw/microb/conf"
)

var Config = conf.GetConf()

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
	var row map[string]interface{}
	var routes []string
	//fmt.Println("RES:", res)
	for res.Next(&row) {
		//fmt.Println("ROW: ", row)
		url := row["uri"].(string)
	    routes = append(routes, url)
	}
	if err != nil {
		fmt.Printf("Rethinkdb: error scanning database results: %s\n", err)
	}
	defer res.Close()
	return routes
	
}

func GetFromDb (url string)  (map[string]interface{}, bool)  {
	session := connectToDb()
	found := false
	db := Config["database"].(string)
	table := Config["table"].(string)
	filters := map[string]interface{}{"uri":url}
	res, err := r.DB(db).Table(table).Filter(filters).Pluck("fields").Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer res.Close()
	//var rows []interface{}
	var rescol map[string]interface{}
	//err = res.All(&rows)
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
		//page_served = rows[0].(map[string]interface{})
		page_served = rescol
	}
	return page_served, found
}

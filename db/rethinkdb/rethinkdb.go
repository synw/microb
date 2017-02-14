package rethinkdb

import (
	"fmt"
	"log"
	 r "gopkg.in/dancannon/gorethink.v2"
	 "github.com/synw/microb/conf"
	 "github.com/synw/microb/datatypes"
	 "github.com/synw/microb/metadata"
)

var Config = conf.GetConf()
var Conn *r.Session
var MainDb *datatypes.Database = metadata.GetMainDatabase()

func init() {
	if MainDb.Type == "rethinkdb" {
		Conn = connectToDb(MainDb)
	}
}

func connectToDb(database *datatypes.Database) (*r.Session) {
	host := database.Host
	port := database.Port
	user := database.User
	pwd := database.Password
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

func GetFromDb(url string)  (map[string]interface{}, bool)  {
	session := Conn
	found := false
	db := Config["domain"].(string)
	filters := map[string]interface{}{"uri":url}
	res, err := r.DB(db).Table("pages").Filter(filters).Pluck("fields").Run(session)
	defer res.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	var rescol map[string]interface{}
	err = res.One(&rescol)
	/*if err == r.ErrEmptyResult {
	    //fmt.Printf("Rethinkdb: no results for uri scan: %s\n", err)
	} else */
	if err != nil && err != r.ErrEmptyResult {
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

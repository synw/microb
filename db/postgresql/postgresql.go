package db

import (
	//"sync"
	//"github.com/synw/microb/conf"
	"github.com/synw/microb/db/datatypes"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

func connectToDb() (host string, db string, user string, pwd string) {
	conn_str := "host="+host+" user="+user+"gorm dbname="+db+" sslmode=disable password="+pwd 
	db, err := gorm.Open("postgres", conn_str)
  	defer db.Close()
  	return db, err
}

func SaveCommand(command *datatypes.Command) {
	return
}

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/synw/microb/conf"
	"github.com/synw/microb/db/datatypes"
	"github.com/synw/microb/utils"
)


var Config = conf.GetConf()
var main_db = conf.GetMainDb()
var Backend = main_db.Type

func Connect(db *datatypes.Database) *gorm.DB {
	conn_str := "host="+db.Host+" user="+db.User+"gorm dbname="+db.Name+" sslmode=disable password="+db.Password 
	pgdb, err := gorm.Open("postgres", conn_str)
  	defer pgdb.Close()
  	if err != nil {
  		fmt.Println(err)
  		utils.PrintEvent("error", "Can't connect to Postgresql")
  	}
  	return pgdb
}

func SaveCommand(command *datatypes.Command) {
	return
}

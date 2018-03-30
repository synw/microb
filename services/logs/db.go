package logs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func initDb(conf *types.Conf) *terr.Trace {
	msgs.Status("Initializing logs database")
	db, err := gorm.Open("sqlite3", conf.LogsDbAddr)
	defer db.Close()
	if err != nil {
		tr := terr.New("services.logs.db.initDb", err)
		return tr
	}
	db.AutoMigrate(&types.Log{})
	return nil
}

func saveToDb(keys []string) {
	for _, _ = range keys {

	}
}

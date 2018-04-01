package logs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var database *gorm.DB

func connect(conf *types.Conf) (*gorm.DB, *terr.Trace) {
	db, err := gorm.Open("sqlite3", conf.LogsDbAddr)
	if err != nil {
		tr := terr.New("services.logs.db.initDb", err)
		return db, tr
	}
	return db, nil
}

func initDb(conf *types.Conf) *terr.Trace {
	msgs.Status("Initializing logs database")
	db, tr := connect(conf)
	if tr != nil {
		tr := terr.Pass("services.logs.db.initDb", tr)
		return tr
	}
	db.AutoMigrate(&types.Log{})
	database = db
	return nil
}

func saveToDb(keys []map[string]interface{}) *terr.Trace {
	for _, key := range keys {
		data := key["data"].(map[string]interface{})
		service := data["service"].(string)
		level := key["level"].(string)
		msg := key["message"].(string)
		class := key["event_class"].(string)
		entry := &types.Log{
			Service: service,
			Level:   level,
			Msg:     msg,
			Class:   class,
		}
		database.Create(entry)
	}
	return nil
}

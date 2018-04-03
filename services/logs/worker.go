package logs

import (
	"encoding/json"
	"github.com/synw/microb/libmicrob/redis"
	"github.com/synw/terr"
	"time"
)

func processLogs(key string) {
	for {
		duration := time.Second * 2
		time.Sleep(duration)
		// get the data from Redis
		keys, err := redis.GetKeys(key)
		if err != nil {
			tr := terr.New("services.logs.worker.processLogs", err)
			tr.Fatal()
		}
		// process the data
		var vals []map[string]interface{}
		for _, key := range keys {
			var data map[string]interface{}
			err := json.Unmarshal(key.([]byte), &data)
			if err != nil {
				tr := terr.New("services.logs.worker.processLogs", err)
				tr.Fatal()
			}
			vals = append(vals, data)
		}
		saveToDb(vals)
	}
}

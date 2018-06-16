package logs

import (
	"encoding/json"
	//"fmt"
	"github.com/synw/microb/redis"
	"github.com/synw/terr"
	"time"
)

func processLogs(key string) {
	for {
		duration := time.Second * 10
		time.Sleep(duration)
		// get the data from Redis
		keys, err := redis.GetKeys(key)
		if err != nil {
			tr := terr.New(err)
			tr.Fatal("Can not get keys in Redis")
		}
		// process the data
		var vals []map[string]interface{}
		for _, key := range keys {
			var data map[string]interface{}
			err := json.Unmarshal(key.([]byte), &data)
			if err != nil {
				tr := terr.New(err)
				tr.Fatal("Can not unmarshal json")
			}
			/*fmt.Println("Log ----------------------")
			for k, v := range data {
				fmt.Println("K", k, v)
			}*/
			class := data["data"].(map[string]interface{})["class"].(string)
			cmdName := data["data"].(map[string]interface{})["command"].(string)
			cmdStatus := data["data"].(map[string]interface{})["command_status"].(string)
			data["class"] = class
			data["command"] = cmdName
			data["command_status"] = cmdStatus
			vals = append(vals, data)
		}
		saveToDb(vals)
	}
}

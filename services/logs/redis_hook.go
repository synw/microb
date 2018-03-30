package logs

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type RedisHook struct {
	Host string
	Port int
	Db   int
}

func newRedisHook() *RedisHook {
	hook := RedisHook{"localhost", 6379, 0}
	return &hook
}

func (hook *RedisHook) Fire(entry *log.Entry) error {
	d := make(map[string]interface{})
	t := entry.Time
	msg := "Log " + " from " + entry.Data["service"].(string) + " " + entry.Message
	level := entry.Level.String()
	d["date"] = t
	d["message"] = msg
	d["level"] = level
	d["data"] = entry.Data
	d["event_class"] = "log"
	data, _ := json.Marshal(d)
	// send to Redis
	key := "log_" + hostname
	err := setKey(key, data)
	if err != nil {
		return err
	}
	return nil
}

func (hook *RedisHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	}
}

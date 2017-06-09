package log

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/synw/centcom"
)

type Hook struct {
	Cli      *centcom.Cli
	Channels map[string]string
}

func NewHook(cli *centcom.Cli, channels map[string]string) *Hook {
	hook := &Hook{
		cli,
		channels,
	}
	return hook
}

func (hook *Hook) Fire(entry *log.Entry) error {
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

	// TODO
	channel := "$channel"

	_, err := hook.Cli.Http.Publish(channel, data)
	if err != nil {
		return err
	}
	return nil
}

func (hook *Hook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	}
}

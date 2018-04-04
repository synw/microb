package logs

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/synw/centcom"
)

type CentHook struct {
	Cli      *centcom.Cli
	Channels map[string]string
}

func NewCentHook(cli *centcom.Cli, channels map[string]string) *CentHook {
	hook := &CentHook{
		cli,
		channels,
	}
	return hook
}

func (hook *CentHook) Fire(entry *log.Entry) error {
	d := make(map[string]interface{})
	t := entry.Time
	msg := entry.Message
	level := entry.Level.String()
	d["date"] = t
	d["message"] = msg
	d["level"] = level
	d["data"] = entry.Data
	d["event_class"] = "log"
	data, _ := json.Marshal(d)

	// TODO
	channel := "$logchan"

	_, err := hook.Cli.Http.Publish(channel, data)
	if err != nil {
		return err
	}
	return nil
}

func (hook *CentHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
	}
}

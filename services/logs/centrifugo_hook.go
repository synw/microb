package logs

import (
	"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/synw/centcom"
)

type CentHook struct {
	Cli      *centcom.Cli
	LogChans map[string]string
	EmitChan string
	Domain   string
}

func NewCentHook(cli *centcom.Cli, logChans map[string]string, channel string, domain string) *CentHook {
	hook := &CentHook{
		cli,
		logChans,
		channel,
		domain,
	}
	return hook
}

func (hook *CentHook) Fire(entry *log.Entry) error {
	if emit == false {
		return nil
	}
	//fmt.Println("EMIT", hook.EmitChan, "//", entry)

	d := make(map[string]interface{})
	t := entry.Time
	msg := entry.Message
	level := entry.Level.String()
	d["Date"] = t
	d["Msg"] = msg
	d["Level"] = level
	d["Data"] = entry.Data
	d["Domain"] = hook.Domain
	data, _ := json.Marshal(d)
	_, err := hook.Cli.Http.Publish(hook.EmitChan, data)
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

package metadata

import (
	//"fmt"
	"github.com/synw/microb/libmicrob/conf"
)

var Config = conf.GetConf()

func GetConfString(param string) string {
	domain := Config[param].(string)
	return domain
}

func GetCommandsBrockers() []string {
	conf_brokers := Config["commands_brokers"].([]interface{})
	nb := len(conf_brokers)
	brokers  := make([]string, nb)
	for i, broker := range(conf_brokers) {
		brokers[i] = broker.(string)
	}
	if len(brokers) == 0 {
		b := []string{}
		return b
	}
	return brokers
}

func IsWebsocketsBrocker() bool {
	brockers := GetCommandsBrockers()
	if brockers != nil {
		for _, b := range(brockers) {
			if b == "websockets" {
				return true
			}
		}
	}
	return false
}

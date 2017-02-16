package encoding

import (
	"time"
	"encoding/json"
	"github.com/synw/microb/libmicrob/datatypes"
)


func DecodeJsonRawMessage(raw *json.RawMessage) (*datatypes.WsMessage, error) {
	var message *datatypes.WsMessage
	byte, err := json.Marshal(raw)
	if err != nil {
		return message, err
	}
	err = json.Unmarshal(byte, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}

func GetCommandFromPayload(message *datatypes.WsMessage) *datatypes.Command {
	data := message.Data
	name := data["name"].(string)
	reason := ""
	if data["reason"].(string) != "" {
		reason = data["reason"].(string)
	}
	from := "websockets"
	now := time.Now()
	status := "pending"
	command := &datatypes.Command{name, from, reason, now, status, nil}
	return command
}

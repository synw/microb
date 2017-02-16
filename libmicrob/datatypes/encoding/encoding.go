package encoding

import (
	"time"
	"encoding/json"
	"github.com/synw/microb/libmicrob/datatypes"
)


func DecodeJsonFeedbackRawMessage(raw *json.RawMessage) (*datatypes.WsFeedbackMessage, error) {
	var message *datatypes.WsFeedbackMessage
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

func DecodeJsonIncomingRawMessage(raw *json.RawMessage) (*datatypes.WsIncomingMessage, error) {
	var message *datatypes.WsIncomingMessage
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

func GetCommandFromPayload(message *datatypes.WsIncomingMessage) *datatypes.Command {
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

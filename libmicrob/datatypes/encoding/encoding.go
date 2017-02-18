package encoding

import (
	"reflect"
	"encoding/json"
	"github.com/ventu-io/go-shortid"
	"github.com/synw/microb/libmicrob/datatypes"
)


func GetType(data interface{}) string {
	t := reflect.TypeOf(data).String()
	return t
}

func GenerateId() string {
	id, _ := shortid.Generate()
	return id
}

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

func MakeWsMsg(command *datatypes.Command) *datatypes.WsMessage {
	data := make(map[string]interface{})
	data["name"] = command.Name
	data["id"] = command.Id
	data["from"] = command.From
	data["reason"] = command.Reason
	data["args"] = command.Args
	ws_msg := &datatypes.WsMessage{"command", "incoming command", data}
	return ws_msg
}

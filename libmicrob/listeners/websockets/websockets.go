package websockets

import (
	//"log"
	"fmt"
	"time"
	"encoding/json"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/centrifugal/gocent"
	"github.com/acmacalister/skittles"
	appevents "github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/datatypes/encoding"
	"github.com/synw/microb/libmicrob/commands"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/state"
)


func credentials() *centrifuge.Credentials {
	secret := state.Server.WebsocketsKey
	user := "microb_"+state.Server.Domain
	timestamp := centrifuge.Timestamp()
	info := ""
	token := auth.GenerateClientToken(secret, user, timestamp, info)
	return &centrifuge.Credentials{
		User:      user,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}
}

func listenForCommands(channel_name string, done chan bool) (centrifuge.Centrifuge, *centrifuge.SubEventHandler) {
	creds := credentials()
	wsURL := "ws://"+state.Server.WebsocketsHost+":"+state.Server.WebsocketsPort+"/connection/websocket"
	
	onMessage := func(sub centrifuge.Sub, rawmsg centrifuge.Message) error {
		if (state.Debug == true) {
			fmt.Println(fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), rawmsg))
		}
		payload, err := encoding.DecodeJsonIncomingRawMessage(rawmsg.Data)
		var msg string
		if err != nil {
			msg = "Error decoding json raw message: "+err.Error()
			appevents.ErrMsg("websockets.listenForCommands()", msg)
		}
		command := commands.GetCommandFromPayload(payload, "websockets")
		msg  = "Command "+skittles.BoldWhite(command.Name)+" received from "+command.From+" via websockets"
		if command.Reason != "" {
			msg = msg+". Reason: "+command.Reason
		}
		appevents.New("event", "websockets", msg)
		commands.Run(command)
		go sendCommandFeedback(command)
		return nil
	}
	/*
	onJoin := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		fmt.Println("JOIN")
		log.Println(fmt.Sprintf("User %s joined channel %s with client ID %s", msg.User, sub.Channel(), msg.Client))
		return nil
	}

	onLeave := func(sub centrifuge.Sub, msg centrifuge.ClientInfo) error {
		log.Println(fmt.Sprintf("User %s with clientID left channel %s with client ID %s", msg.User, msg.Client, sub.Channel()))
		return nil
	}
	*/
	onPrivateSub := func(c centrifuge.Centrifuge, req *centrifuge.PrivateRequest) (*centrifuge.PrivateSign, error) {
		info := ""
		sign := auth.GenerateChannelSign(state.Server.WebsocketsKey, req.ClientID, req.Channel, info)
		privateSign := &centrifuge.PrivateSign{Sign: sign, Info: info}
		return privateSign, nil
	}

	events := &centrifuge.EventHandler{
		OnPrivateSub: onPrivateSub,
	}
	
	subevents := &centrifuge.SubEventHandler{
		OnMessage: onMessage,
		/*OnJoin:    onJoin,
		OnLeave:   onLeave,*/
	}
	c := centrifuge.NewCentrifuge(wsURL, creds, events, centrifuge.DefaultConfig)
	return c, subevents
}

func sendCommandFeedback(command *datatypes.Command) {
	secret := state.Server.WebsocketsKey
	host := state.Server.WebsocketsHost
	port := state.Server.WebsocketsPort
	purl := fmt.Sprintf("http://%s:%s", host, port)
	// connect to Centrifugo
	client := gocent.NewClient(purl, secret, 5*time.Second)
	var errstr string
	if command.Error == nil {
		errstr = ""
	} else {
		errstr = command.Error.Error()
	}
	data := make(map[string]interface{})
	data["command"] = command.Name
	data["id"] = command.Id
	data["from"] = command.From
	data["reason"] = command.Reason
	if len(command.ReturnValues)  > 0 {
		data["return_values"] = command.ReturnValues
	}
	eventstr := &datatypes.WsFeedbackMessage{"command_feedback", command.Status, errstr, data}
	event, err := json.Marshal(eventstr)
	channel := "$"+state.Server.Domain+"_feedback"
	if (state.Debug == true) {
		fmt.Println("Sending feedback message to", channel, eventstr, "\nDATA:", data)
	}
	_, err = client.Publish(channel, event)
	if err != nil {
		appevents.Error("websockets.SendCommandFeedback()", err)
	}
	return
}

func ListenToIncomingCommands(channel_name string) {
	var done chan bool
	// connect to channel
	c, subevents := listenForCommands(channel_name, done)
	defer c.Close()
	err := c.Connect()
	if err != nil {
		appevents.New("error", "listeners.websockets.Listen()", err.Error())
	}
	// suscribe to channel
	_, err = c.Subscribe(channel_name, subevents)
	if err != nil {
		appevents.New("error", "listeners.websockets.Listen()", err.Error())
	}
	// sit here and wait
	<- done
}

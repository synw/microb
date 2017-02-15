package websockets

import (
	"log"
	"fmt"
	"encoding/json"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/synw/microb/libmicrob/conf"
	appevents "github.com/synw/microb/libmicrob/events"
)


var Config = conf.GetConf()
var SecretKey string = Config["centrifugo_secret_key"].(string)

func credentials() *centrifuge.Credentials {
	secret := SecretKey
	user := "microb_"+Config["domain"].(string)
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

type WsMsg struct {
	EventClass string `json:"event_class"`
	Data map[string]interface{} `json:"data"`
}

type Model struct {
	EventClass  json.RawMessage `json:"event_class"`
	Data string          `json:"data"`
}

func (m *Model) Name() string {
	return string(m.EventClass)
}

func listenForCommands(channel_name string, done chan bool) (centrifuge.Centrifuge, *centrifuge.SubEventHandler) {
	creds := credentials()
	wsURL := "ws://"+Config["centrifugo_host"].(string)+":"+Config["centrifugo_port"].(string)+"/connection/websocket"
	
	onMessage := func(sub centrifuge.Sub, msg centrifuge.Message) error {
		log.Println(fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), msg))
		//x := new(WsMsg)
		
		/*
		var s string
		d := json.Marshal(msg.Data)
	    err := json.Unmarshal(d, &s)
	    if err == nil {
	        x.EventClass = s
	    }
		fmt.Println(x)
		*/
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
		sign := auth.GenerateChannelSign(SecretKey, req.ClientID, req.Channel, info)
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

func ListenToIncomingCommands(channel_name string) {
	var done chan bool
	c, subevents := listenForCommands(channel_name, done)
	defer c.Close()
	err := c.Connect()
	if err != nil {
		appevents.New("error", "listeners.websockets.Listen()", err.Error())
	}
	_, err = c.Subscribe(channel_name, subevents)
	if err != nil {
		appevents.New("error", "listeners.websockets.Listen()", err.Error())
	}
	<- done
}

package datatypes

import (
	"time"
	"strconv"
	//"errors"
	"encoding/json"
	format "github.com/synw/microb/libmicrob/events/format/methods"
)


type Server struct {
	Domain string
	Host string
	Port string
	WebsocketsHost string
	WebsocketsPort string
	WebsocketsKey string
	PagesDb *Database
	HitsDb *Database
	CommandsDb *Database
	Runing bool
}

func (s Server) Format() string {
	var msg string
	msg = msg+" - Domain: "+s.Domain+"\n"
	msg = msg+" - Host: "+s.Host+"\n"
	msg = msg+" - Port: "+s.Port+"\n"
	msg = msg+" - Websockets host: "+s.Port+"\n"
	msg = msg+" - Websockets port: "+s.Port+"\n"
	msg = msg+" - Pages database: "+s.PagesDb.Name+"\n"
	msg = msg+" - Hits database: "+s.HitsDb.Name+"\n"
	msg = msg+" - Commands database: "+s.CommandsDb.Name+"\n"
	if (s.Runing) {
		msg = msg+" * Server is runing"
	} else {
		msg = msg+" * Server is not runing"
	}
	return msg
}

type Database struct {
	Type string
	Name string
	Host string
	Port string
	User string
	Password string
	Roles []string
}

type Page struct {
    Url string
    Title string
    Content  string
}

func (p Page) Format() string {
	var msg string
	msg = msg+" - Url: "+p.Url+"\n"
	msg = msg+" - Title: "+p.Title+"\n"
	msg = msg+" - Content: "+p.Content
	return msg
}

type Event struct {
	Class string
	From string
	Message string
	Data map[string]interface{}
}

type Command struct {
	Id string
	Name string
	From string
	Reason string
	Date time.Time
	Args []interface{}
	Status string
	Error error
	ReturnValues []interface{}
}

type WsMessage struct { 
	EventClass string `json:"event_class"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

type WsIncomingMessage struct {
	RawMessage  json.RawMessage 
	EventClass string `json:"event_class"`
	Message string `json:"message"`
	Data map[string]interface{} `json:"data"`
}

type WsFeedbackMessage struct {
	EventClass string `json:"event_class"`
	Status string `json:"status"`
	Error string `json:"error"`
	Data map[string]interface{} `json:"data"`
}

type HttpResponse struct {
	*Server
	Url string
	Content string
	Size int
	StatusCode int
	Duration time.Duration
}

func (r HttpResponse) Format() string {
	msg := r.Content+"\n"
	msg = msg+format.FormatStatusCode(r.StatusCode)+" "+r.Url+"\n"
	h := float64(r.Size)/1000
	f := strconv.FormatFloat(h, 'f', 2, 64)
	msg = msg+"Size: "+f+" Ko - "
	d := r.Duration.String()
	msg = msg+d
	return msg
}

type HttpRequest struct {
	*Server
	*Page
	StatusCode int
	//Error errors.Err
}

type HttpRequestMetric struct {
	*Server
	Url string
	ProcessingTime int
	TransportTime int
	TotalTime int
	StatusCode int
}

func (m HttpRequestMetric) Format() string {
	msg := format.FormatStatusCode(m.StatusCode)+" "+m.Url+"\n"
	msg = msg+" - Server processing time: "+strconv.Itoa(m.ProcessingTime)+" ms"+"\n"
	msg = msg+" - Transport time: "+strconv.Itoa(m.TransportTime)+" ms\n"
	msg = msg+" - Total time: "+strconv.Itoa(m.TotalTime)+" ms"
	return msg
}

type HttpStressReport struct {
	*Server
	NumRequests int
	Size int
	Duration time.Duration
}

func (r HttpStressReport) Format() string {
	msg := "Stressing server "+r.Domain+" with "+strconv.Itoa(r.NumRequests)+" requests\n"
	h := float64(r.Size)/1000
	f := strconv.FormatFloat(h, 'f', 2, 64)
	msg = msg+" - Total size: "+f+" Ko\n"
	d := r.Duration.String()
	msg = msg+" - Total duration: "+d
	return msg
}

/*
type GormPage struct {
	Domain string
	Url string `gorm:"not null;unique"`
    Title string
    Content  string
	
}

type Hit struct {
	Url string
	Method string
	Ip string
	User_agent string
	Referer string
	Date time.Time
}

type DataChanges struct {
	Msg string
	Type string
	Values map[string]interface{}
}*/

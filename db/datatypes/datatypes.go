package datatypes

import (
	"time"
)

type DataChanges struct {
	Msg string
	Type string
	Values map[string]interface{}
}

type Command struct {
	Name string
}

type Hit struct {
	Url string
	Method string
	Ip string
	User_agent string
	Referer string
	Date time.Time
}

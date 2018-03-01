package types

import (
	"github.com/synw/terr"
	"time"
)

type Service struct {
	Name     string
	Cmds     []Cmd
	Init     func(bool, int, bool) *terr.Trace
	Dispatch func(*Cmd) *Cmd
}

type Cmd struct {
	Id           string
	Name         string
	Date         time.Time
	Args         []interface{}
	Status       string
	Trace        *terr.Trace
	ReturnValues []interface{}
}

type Event struct {
	Id      string
	Class   string
	Date    time.Time
	Msg     string
	Service *Service
	Cmd     *Cmd
	Trace   *terr.Trace
	Data    map[string]interface{}
}

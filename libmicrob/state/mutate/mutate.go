package mutate

import (
	"strconv"
	"errors"
    "github.com/synw/microb/libmicrob/state"
    "github.com/synw/microb/libmicrob/events"
)


func Verbosity(lvl string) string {
	v, _ := strconv.Atoi(lvl)
	state.Verbosity = v
	msg := "Verbosity is set to "+lvl
	events.State("mutate.Verbosity", msg)
	return msg
}

func Debug(lvl string) (string, error) {
	var msg string
	if lvl == "true" {
		state.Debug = true
		msg = "Debug is set to "+lvl
		events.State("mutate.Debug", msg)
	} else {
		msg = "Invalid value for set debug: "+lvl
		events.ErrMsg("mutate.Debug", msg)
		err := errors.New(msg)
		return msg, err
	}
	return msg, nil
}

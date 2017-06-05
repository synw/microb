package services

import (
	"errors"
	"fmt"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/terr"
	"reflect"
)

func Dispatch(cmd *datatypes.Command, c chan *datatypes.Command) {
	com := &datatypes.Command{}
	for n, _ := range initDispatch {
		if cmd.Service == n {
			cm, _ := Call(initDispatch, n, cmd)
			com = cm[0].Interface().(*datatypes.Command)
		}
	}
	_ = events.Cmd(cmd)
	c <- com
	return
}

func InitServices(dev bool, verbosity int) *terr.Trace {
	if state.Verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	// declare services
	for _, el := range state.Conf["services"].([]interface{}) {
		name := el.(string)
		s := New(name)
		state.Services = append(state.Services, s)
		if state.Verbosity > 1 {
			fmt.Println("Registering service", name)
		}
		// init service state
		for n, _ := range initState {
			if name == n {
				// initilize service state
				_, tr := Call(initState, n, dev, verbosity)
				if tr != nil {
					tr = terr.Add("initServices", errors.New("Unable to initialize "+name+" service"))
					return tr
				}
			}
		}
		if state.Verbosity > 0 {
			msg := name + " service initialized"
			events.Msg("info", "services.Init", msg)
		}
	}
	return nil
}

func New(name string, deps ...[]*datatypes.Service) *datatypes.Service {
	req := []*datatypes.Service{}
	if len(deps) > 0 {
		for _, dep := range deps[0] {
			req = append(req, dep)
		}
	}
	s := &datatypes.Service{name, req}
	return s
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, tr *terr.Trace) {
	f := reflect.ValueOf(m[name])
	var v []reflect.Value
	if len(params) != f.Type().NumIn() {
		err := errors.New("The number of params is not adapted.")
		tr := terr.New("services.Call", err)
		return v, tr
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return in, nil
}

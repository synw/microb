package services

import (
	"errors"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/state"
	"github.com/synw/terr"
	"reflect"
)

var ValidCommands []string

func Dispatch(cmd *datatypes.Command, c chan *datatypes.Command) {
	com := &datatypes.Command{}
	found := false
	for n, _ := range initDispatch {
		if cmd.Service == n {
			cm, _ := Call(initDispatch, n, cmd)
			com = cm[0].Interface().(*datatypes.Command)
			found = true
			break
		}
	}
	if !found {
		com.Status = "error"
		err := errors.New("Can not find service " + cmd.Service)
		tr := terr.New("services.Dispatch", err)
		com.Trace = tr
	}
	_ = events.Cmd(cmd)
	c <- com
	return
}

func InitServices(dev bool, verbosity int) *terr.Trace {
	if state.Verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	servs := state.Conf["services"].([]interface{})
	servs = append(servs, "info")
	// declare services
	for _, el := range servs {
		// register service
		name := el.(string)
		state.Services[name] = All[name]
		// register service commands
		ValidCommands = append(ValidCommands, All[name].Cmds...)
		// init service state
		for n, _ := range initState {
			if name == n {
				// initialize service state
				_, tr := Call(initState, n, dev, verbosity)
				if tr != nil {
					tr = terr.Add("initServices", errors.New("Unable to initialize "+name+" service"))
					return tr
				}
			}
		}
		if state.Verbosity > 0 {
			if name != "info" {
				msg := color.BoldWhite(name) + " service initialized"
				events.Msg("info", "services.Init", msg)
			}
		}
	}
	return nil
}

func New(name string, cmds []string, deps ...[]*datatypes.Service) *datatypes.Service {
	req := []*datatypes.Service{}
	if len(deps) > 0 {
		for _, dep := range deps[0] {
			req = append(req, dep)
		}
	}
	s := &datatypes.Service{name, cmds, req}
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

// used by cli
func getAllValidCommands() []string {
	vc := []string{}
	for _, s := range All {
		vc = append(vc, All[s.Name].Cmds...)
	}
	return vc
}

func CmdIsValid(cmd *datatypes.Command) bool {
	name := cmd.Name
	for _, c := range getAllValidCommands() {
		if c == name {
			return true
		}
	}
	return false
}

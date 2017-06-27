package services

import (
	"errors"
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var ValidCommands []string
var isDev = false

func GetService(name string) *types.Service {
	serv := AllDev[name]
	return serv
}

func InitServices(dev bool, verbosity int, servs []string) (map[string]*types.Service, []*terr.Trace) {
	if verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	services := make(map[string]*types.Service)
	var trs []*terr.Trace
	for _, name := range servs {
		service := GetService(name)
		tr := service.Init(dev, verbosity)
		if tr != nil {
			tr = terr.Add("initServices", errors.New("Unable to initialize "+name+" service"))
			msg := tr.Formatc()
			events.Err(name, "services.InitServices", msg, tr.ToErr())
			trs = append(trs, tr)
		} else {
			msg := color.BoldWhite(name) + " service initialized"
			events.State(name, "services.Init", msg, nil)
		}
		services[name] = service
		ValidCommands = append(ValidCommands, service.Cmds...)
	}
	return services, trs
}

func Dispatch(cmd *types.Command, c chan *types.Command) {
	service := GetService(cmd.Service)
	com := service.Dispatch(cmd)
	c <- com
}

/*
func Dispatch(cmd *types.Command, c chan *types.Command) {
	com := &types.Command{}
	found := false
	di := dispatch
	if isDev == true {
		di = dispatchDev
	}
	for n, _ := range di {
		if cmd.Service == n {
			cm, _ := Call(di, n, cmd)
			com = cm[0].Interface().(*types.Command)
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

	if cmd.Service == "http" {
		com = http.Service.Dispatch(cmd)
	}

	c <- com
	return
}

/*
func InitServices(dev bool, verbosity int, servs []string) (map[string]*types.Service, []*terr.Trace) {
	isDev = dev
	var trs []*terr.Trace
	services := make(map[string]*types.Service)
	if verbosity > 0 {
		fmt.Println("Initializing services ...")
	}
	servs = append(servs, "info")
	// declare services
	for _, name := range servs {
		// register service
		s := All[name]
		if dev == true {
			s = AllDev[name]
		}
		services[name] = s
		// register service commands
		ValidCommands = append(ValidCommands, s.Cmds...)
		// init service state
		st := initState
		if dev == true {
			st = initStateDev
		}
		for n, _ := range st {
			if name == n {
				// initialize service state
				_, tr := Call(st, n, dev, verbosity)
				if tr != nil {
					tr = terr.Add("initServices", errors.New("Unable to initialize "+name+" service"))
					msg := tr.Formatc()
					events.Err(name, "services.InitServices", msg, tr.ToErr())
					trs = append(trs, tr)
				} else {
					msg := color.BoldWhite(name) + " service initialized"
					events.State(name, "services.Init", msg, nil)
				}
			}
		}

		http.Service.Init(true, 2)

	}
	return services, trs
}*/
/*
func New(name string, cmds []string, deps ...[]*types.Service) *types.Service {
	req := []*types.Service{}
	if len(deps) > 0 {
		for _, dep := range deps[0] {
			req = append(req, dep)
		}
	}
	s := &types.Service{name, cmds, req}
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

	for _, el := range result {
		t := reflect.TypeOf(el.Interface())
		if fmt.Sprintf("%s", t) == "*terr.Trace" {
			if el.Elem().IsValid() == true {
				err := errors.New("Unable to initilaize service " + name)
				tr := terr.Add("services.Call", err)
				return in, tr
			}
		}
	}
	return in, nil
}*/

// used by cli
func getAllValidCommands() []string {
	vc := []string{}
	for _, s := range All {
		vc = append(vc, All[s.Name].Cmds...)
	}
	return vc
}

func CmdIsValid(cmd *types.Command) bool {
	name := cmd.Name
	for _, c := range getAllValidCommands() {
		if c == name {
			return true
		}
	}
	return false
}

package cmd

import (
	//sid "github.com/SKAhack/go-shortid"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func New(name string, service *types.Service, from string, args ...interface{}) *types.Cmd {
	// TO FIX
	//id, _ := sid.Generate()
	id := "id"
	date := time.Now()
	status := "pending"
	var tr *terr.Trace
	var rvs []interface{}
	var exec func(args ...interface{})
	cmd := &types.Cmd{
		id,
		name,
		date,
		service,
		args,
		from,
		status,
		tr,
		rvs,
		exec,
	}
	return cmd
}

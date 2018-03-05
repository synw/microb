package libmicrob

import (
	//sid "github.com/SKAhack/go-shortid"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func New(name string, service string, from string, args ...interface{}) *types.Cmd {
	// TO FIX
	//id, _ := sid.Generate()
	id := "id"
	date := time.Now()
	status := "pending"
	var tr *terr.Trace
	var rvs []interface{}
	cmd := &types.Cmd{
		id,
		name,
		date,
		args,
		from,
		status,
		tr,
		rvs,
	}
	return cmd
}

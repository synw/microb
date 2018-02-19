# Microb

Platform to manage observable services. A service can receive and execute commands from the 
[Microb terminal client](https://github.com/synw/microb-cli)

#### Dependencies

- Websockets server: [Centrifugo](https://github.com/centrifugal/centrifugo)

#### Terminal client

A [terminal client](https://github.com/synw/microb-cli) is used to control Microb instances

## Install

   ```bash
   go get github.com/synw/microb
   go get github.com/synw/microb-cli
   go get github.com/synw/terr
   ```
   
[Install Centrifugo](https://fzambia.gitbooks.io/centrifugal/content/server/start.html)

## Usage

To write a new service you only need a `manifest/manifest.go` inside your project to declare the Microb service: 
ex: an http service:

   ```go
package manifest

import (
	   "github.com/synw/microb-http/cmd"
	   "github.com/synw/microb-http/state"
	   "github.com/synw/microb/libmicrob/types"
	   "github.com/synw/terr"
)

var Service *types.Service = &types.Service{
	   "http", // name of the service
	   []string{"start", "stop", "parse_templates"}, //commands
	   ini, // function to initialize
	   dispatch, //function to dispatch commands
}

func ini(dev bool, verbosity int, start bool) *terr.Trace {
	   return state.Init(dev, verbosity, start)
}

func dispatch(c *types.Command) *types.Command {
	   return cmd.Dispatch(c)
}
   ```
   
Then you can write commands that will be run from the client using a `cmds` package: ex: a simple ping command in 
`cmds/infos/info.go`:

   ```go
package info

import (
	   "github.com/synw/microb/libmicrob/types"
)

func Dispatch(cmd *types.Command) *types.Command {
	   com := &types.Command{}
	   if cmd.Name == "ping" {
	   	return Ping(cmd)
	   }
	   return com
}

func Ping(cmd *types.Command) *types.Command {
	   var resp []interface{}
	   resp = append(resp, "PONG")
	   cmd.ReturnValues = resp
	   cmd.Status = "success"
	   return cmd
}
   ```

Check the [Microb http service](https://github.com/synw/microb-http)

Then register your service in `microb/services/manifest.go`:

   ```go
package services

import (
	   http "github.com/synw/microb-http/manifest"
	   "github.com/synw/microb/libmicrob/types"
	   "github.com/synw/microb/services/info"
)

var All = map[string]*types.Service{
	   "info": info.Service,
	   "http": http.Service,
}
   ```
   
Declare it in `microb/config.json` to enable it:

   ```javascript
{
	   "centrifugo_addr":"localhost:8001",
	   "centrifugo_key":"secret_key",
	   "verbosity": 1,
	   "name": "localhost",
	   "services": ["http"]
}
   ```
   
Compile Microb
   
## Run

   ```
   ./microb -d -s -v 1
   ```

Flags:

- `d`: for dev mode. It will use `dev_config.json`
- `s`: to start the services instead of doing it from the client
- `v`: verbosity level

#### External libraries

- [Viper](https://github.com/spf13/viper): configuration management
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers
- [Go-short-id](https://github.com/ventu-io/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors
- [Chi](https://github.com/pressly/chi): http router

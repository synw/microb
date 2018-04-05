# Microb

A microbiota where to deploy observable Go services. A service can receive and execute commands from the 
[Microb terminal client](https://github.com/synw/microb-cli)

#### Dependencies

- [Centrifugo](https://github.com/centrifugal/centrifugo): the websockets server
- [Redis](https://redis.io/): the famous in memory key/value store

#### Terminal client

A [terminal client](https://github.com/synw/microb-cli) is used to control Microb instances

## Install

   ```bash
   go get github.com/synw/microb
   ```
   
[Install Centrifugo](https://fzambia.gitbooks.io/centrifugal/content/server/start.html)

## Usage

Create a Go package for your service and put a `manifest` folder in it. First file is for service
initilization: create a `init.go` file: example for a hello world service:

   ```go
   package hello_world

   import (
      "github.com/synw/microb/libmicrob/types"
   )

   var Service *types.Service = &types.Service{
	  "hello_world", // service name
	   getCmds(), // function to get the service commands
	   initService, // function to initialize the service
   }
   
   func initService(dev bool, start bool) error {
       // this service does not need anything
       // special for its initialization
	   return nil
   }
   ```
   
Make a `cmds.go` file to define your service commands:

   ```go
   package hello_world

   import (
	   "github.com/synw/microb/libmicrob/types"
   )

   func getCmds() map[string]*types.Cmd {
	   cmds := make(map[string]*types.Cmd)
	   cmds["hello"] = hello()
	   return cmds
   }

   func hello() *types.Cmd {
       // define the command and attach its running function
	   cmd := &types.Cmd{Name: "hello", Exec: runHello}
	   return cmd
   }
   
   func runHello(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
       // this function will be run on command call
	   var resp []interface{}
	   resp = append(resp, "Hello world")
	   // the command will return "Hello world"
	   cmd.ReturnValues = resp
	   cmd.Status = "success"
	   c <- cmd
   }
   ```
   
Now declare the service client-side: `go get github.com/synw/microb-cli` and open
`services/manifest.go` to add your service:

   ```go
   package services

   import (
       "github.com/synw/microb/libmicrob/types"
	   http "github.com/synw/microb-http/manifest"
	   infos "github.com/synw/microb/services/infos"
	   hello "github.com/me/hello_world/manifest"
   )

   var services = map[string]*types.Service{
      // other services
      "infos": infos.Service,
      "http":  http.Service,
      // declare my service
      "hello_world":  hello.Service,
   }
   ```
   
Compile the client and the server and it's ready to use: sending the command `hello` from
the client will return "Hello world"

## Available sevices

- [Http](https://github.com/synw/microb-http): an http server
- [Mail](https://github.com/synw/microb-mail): send mails

#### External libraries

- [Viper](https://github.com/spf13/viper): configuration management
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Redigo](https://github.com/garyburd/redigo): drivers for Redis
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers
- [Fsm](https://github.com/looplab/fsm): finite state machine lib
- [Go-short-id](https://github.com/ventu-io/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors

# Microb

A microbiota where to deploy observable Go services. A service can receive and execute commands from the 
[Microb terminal client](https://github.com/synw/microb-cli). Features:

- **Remote commands**: all services can have remote commands to talk to the Microb server
- **Logs**: all logs are stored in an sqlite database

#### Dependencies

- [Centrifugo](https://github.com/centrifugal/centrifugo): the websockets server
- [Redis](https://redis.io/): the famous in memory key/value store

#### Terminal client

A [terminal client](https://github.com/synw/microb-cli) is used to control Microb instances

## Install

   ```bash
   go get github.com/synw/microb
   ```
   
[Install Centrifugo](https://fzambia.gitbooks.io/centrifugal/content/server/start.html) and 
Redis

## Example usage

Let's do a hello world service. Create a `hello_world` Go package for your service and
a `manifest` folder in it.

### 1. Initialize the service

Create an `init.go` file in the manifest folder:

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

### 2. Create some commands
   
Create a `cmds.go` file in the manifest folder to define your service commands:

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
   
   func runHello(cmd *types.Cmd, c chan *types.Cmd) {
       // this function will be run on command call
	   var resp []interface{}
	   resp = append(resp, "Hello world")
	   // the command will return "Hello world"
	   cmd.ReturnValues = resp
	   cmd.Status = "success"
	   c <- cmd
   }
   ```

### 3. Declare the service server side
   
Declare the service in Microb: open `services/manifest.go` in the Microb package and
add your service:
   
   ```go
   package services

   import (
       "github.com/synw/microb/libmicrob/types"
	   http "github.com/synw/microb-http/manifest"
	   infos "github.com/synw/microb/services/infos"
	   hello_world "github.com/me/hello_world/manifest"
   )

   var Services = map[string]*types.Service{
	   "infos": infos.Service,
	   "http":  http.Service,
	   "hello_world":  hello_world.Service,
   }
   ```

### 4. Declare the service client side

Now declare the service client-side: `go get github.com/synw/microb-cli` and open
`services/manifest.go` to add your service in the same way:

   ```go
   package services

   import (
       "github.com/synw/microb/libmicrob/types"
	   http "github.com/synw/microb-http/manifest"
	   infos "github.com/synw/microb/services/infos"
	   hello_world "github.com/me/hello_world/manifest"
   )

   var services = map[string]*types.Service{
      // other services
      "infos": infos.Service,
      "http":  http.Service,
      // declare my service
      "hello_world":  hello_world.Service,
   }
   ```

### 5. Update the config files
   
Update the client and server `config.json` to enable your service:

   ```
   {
   ...
   "services":["infos", "http", "hello_world"]
   }
   ```

### Final step: compilation
   
Compile the client and the server and it's ready to use: sending the command `hello` from
the client will return "Hello world"

## Available sevices

- [Http](https://github.com/synw/microb-http): an http server
- [Mail](https://github.com/synw/microb-mail): send mails

## Screenshot

The upper terminal shows the server output and the lower one shows the client output:

![Screenshot](https://raw.githubusercontent.com/synw/microb/master/docs/img/screenshot.png)

#### External libraries

- [Viper](https://github.com/spf13/viper): configuration management
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Redigo](https://github.com/gomodule/redigo): drivers for Redis
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers
- [Fsm](https://github.com/looplab/fsm): finite state machine lib
- [Go-short-id](https://github.com/SKAhack/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors

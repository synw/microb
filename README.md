# Microb

A microbiota where to deploy observable Go services. A service can receive and execute commands from the 
[Microb terminal client](https://github.com/synw/microb-cli)

**Note**: a full rewrite is in progress: see the v3 branch for a working version

#### Dependencies

- Websockets server: [Centrifugo](https://github.com/centrifugal/centrifugo)

#### Terminal client

A [terminal client](https://github.com/synw/microb-cli) is used to control Microb instances

## Install

   ```bash
   go get github.com/synw/microb
   ```
   
[Install Centrifugo](https://fzambia.gitbooks.io/centrifugal/content/server/start.html)


#### External libraries

- [Viper](https://github.com/spf13/viper): configuration management
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers
- [Fsm](https://github.com/looplab/fsm): finite state machine lib
- [Go-short-id](https://github.com/ventu-io/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors

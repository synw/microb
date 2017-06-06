# Microb

Toolbox to manage observable services

#### Dependencies

- Websockets server: [Centrifugo](https://github.com/centrifugal/centrifugo)

#### Terminal client

A [terminal client](https://github.com/synw/microb-cli) is used to control Microb instances

## Install

   ```bash
   go get github.com/synw/microb
   go get github.com/synw/microb-cli
   ```
   
[Install Centrifugo](https://fzambia.gitbooks.io/centrifugal/content/server/start.html)

## Usage

To write a new service use this folder structure under `services/yourservice`:

- `datatypes`: the types used in the service
- `commands`: all the commands used to mutate or observe the service
- `conf`: the service's configuration
- `state`: manage the service's state

Check the http service example at `services/httpServer`

Then register your service in `services/manifest.go` and compile. Write the cli commands and compile the client.

#### External libraries

- [Viper](https://github.com/spf13/viper): configuration management
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers
- [Go-short-id](https://github.com/ventu-io/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors
- [Chi](https://github.com/pressly/chi): http router

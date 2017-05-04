Microb
======

An API server that is:

- **Observable**: the server can stream runtime info, and send introspections results

- **Mutable**: ability to change the server's state at runtime. A [terminal client](https://github.com/synw/microb-cli)
is available

- **Decoupled**: the databases, the static files server, the websockets server and the api server 
can be located anywhere. 

Dependencies
------------

- Websockets server: [Centrifugo](https://github.com/centrifugal/centrifugo)

Supported databases
-------------------

- Rethinkdb
- Postgresql [WIP]
- Sqlite [WIP]


Terminal client
---------------

A [terminal client](https://github.com/synw/microb-cli) is available to control Microb servers

User interface
--------------

A [user interface](https://github.com/synw/microb-ui) using the Caddy static file server is available

Data flow
---------

![Microb data flow](https://raw.github.com/synw/django-microb/master/docs/img/microb.png)

Note: all components are independant and can be located at any machine: 
database, api server, static server, websockets server and terminal client

External libraries used
-----------------------

- [Chi](https://github.com/pressly/chi): http router
- [Viper](https://github.com/spf13/viper): configuration management
- [Gorm](https://github.com/jinzhu/gorm): orm for relational databases
- [Gorethink](https://github.com/GoRethink/gorethink): Rethinkdb drivers
- [Go-short-id](https://github.com/ventu-io/go-shortid): unique ids generation
- [Skittles](https://godoc.org/github.com/acmacalister/skittles): terminal colors
- [Centrifuge-go](https://github.com/centrifugal/centrifuge-go): Centrifugo server side drivers
- [Gocent](https://github.com/centrifugal/gocent): Centrifugo client side drivers

Microb
======

An API server that is:

- **Observable**: the server can stream runtime info, and send introspections results

- **Mutable**: ability to change the server's state at runtime. A [terminal client](https://github.com/synw/microb-cli)
is available

- **Decoupled**: the databases, the static files server, the websockets server and the Api server 
can be located anywhere. 

External servers used
---------------------

- Static files: [Caddy](https://github.com/mholt/caddy)

- Websockets: [Centrifugo](https://github.com/centrifugal/centrifugo)

Supported databases
-------------------

- Rethinkdb

- Postgresql

- Sqlite

package datatypes

import (
	t "github.com/synw/microb/libmicrob/datatypes"
	"net/http"
)

type HttpServer struct {
	Domain   string
	Server   *t.Server
	Host     string
	Port     int
	Instance *http.Server
	Running  bool
	Cors     string
}

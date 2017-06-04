package datatypes

import (
	"net/http"
)

type HttpServer struct {
	Domain   string
	Host     string
	Port     int
	Instance *http.Server
	Running  bool
	Cors     string
}

type Document struct {
	Url  string
	Data string
}

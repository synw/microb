package services

import (
	"github.com/synw/microb/libmicrob/datatypes"
)

func New(name string, deps ...[]*datatypes.Service) *datatypes.Service {
	req := []*datatypes.Service{}
	if len(deps) > 0 {
		for _, dep := range deps[0] {
			req = append(req, dep)
		}
	}
	s := &datatypes.Service{name, req}
	return s
}

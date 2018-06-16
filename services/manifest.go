package services

import (
	//feeds "github.com/synw/microb-feeds/manifest"
	http "github.com/synw/microb-http/manifest"
	mail "github.com/synw/microb-mail/manifest"
	"github.com/synw/microb/types"
	infos "github.com/synw/microb/services/infos"
	logs "github.com/synw/microb/services/logs"
)

var Services = map[string]*types.Service{
	"logs":  logs.Service,
	"infos": infos.Service,
	"http":  http.Service,
	"mail":  mail.Service,
	//"feeds": feeds.Service,
}

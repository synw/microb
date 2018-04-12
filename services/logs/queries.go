package logs

import (
	"github.com/synw/microb/libmicrob/types"
)

func getLogs(limit int) []types.Log {
	logData := []types.Log{}
	database.Limit(limit).Order("created_at desc").Find(&logData)
	return logData
}

package logs

import (
	"github.com/synw/microb/libmicrob/types"
)

func getLogs(limit int) []types.Log {
	logData := []types.Log{}
	database.Limit(limit).Order("created_at desc").Find(&logData)
	return logData
}

func getErrs(limit int) []types.Log {
	logData := []types.Log{}
	database.Where("level = ?", "error").Order("created_at desc").Limit(limit).Find(&logData)
	return logData
}

func getWarns(limit int) []types.Log {
	logData := []types.Log{}
	database.Where("level = ?", "warning").Order("created_at desc").Limit(limit).Find(&logData)
	return logData
}

func getState(limit int) []types.Log {
	logData := []types.Log{}
	database.Where("class = ?", "state").Order("created_at desc").Limit(limit).Find(&logData)
	return logData
}

func getCommandsFromDb(limit int) []types.Log {
	logData := []types.Log{}
	database.Where("class in (?)", []string{"command_in", "command_out"}).Order("created_at desc").Limit(limit).Find(&logData)
	return logData
}

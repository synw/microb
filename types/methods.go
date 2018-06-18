package types

import (
	"fmt"
	"strings"
)

func (event *Event) Print() {
	var lines []string
	lines = append(lines, "Event:")
	lines = append(lines, "- Service: "+event.Service)
	lines = append(lines, "- Class: "+event.Class)
	lines = append(lines, "- Message: "+event.Msg)
	lines = append(lines, "- Date: "+event.Date.Format("2006-01-02 15:04:05"))
	if event.Cmd != nil {
		lines = append(lines, "- Command: "+event.Cmd.Name)
	} else {
		lines = append(lines, "- Not a command")
	}
	lines = append(lines, "- Log level: "+event.LogLvl)
	if len(event.Trace.Errors) > 0 {
		lines = append(lines, "- Trace: \n"+event.Trace.Msg())
	} else {
		lines = append(lines, "- No errors")
	}
	lines = append(lines, "- Data: ")
	msg := strings.Join(lines, "\n")
	fmt.Println(msg, event.Data)

}

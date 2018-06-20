package msgs

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/types"
	//"github.com/synw/terr"
	"strings"
)

func Warning(txt string, output ...string) {
	msg := "[" + color.Magenta("Warning") + "] " + txt
	fmt.Println(msg)
}

func Status(txt string, output ...string) {
	msg := "[" + color.Blue("Status") + "] " + txt
	fmt.Println(msg)
}

func State(txt string, output ...string) {
	msg := "[" + color.Yellow("State") + "] " + txt
	fmt.Println(msg)
}

func Ready(txt string, output ...string) {
	msg := "[" + color.BoldGreen("Ready") + "] " + txt
	fmt.Println(msg)
}

func Msg(txt string, output ...string) {
	fmt.Println(txt)
}

func Ok(txt string, output ...string) {
	msg := "[" + color.Green("Ok") + "] " + txt
	fmt.Println(msg)
}

func Timeout(txt string, output ...string) {
	msg := "[" + color.BoldRed("Timeout") + "] " + txt
	fmt.Println(msg)
}

func Error(txt string, output ...string) {
	msg := "[" + color.BoldRed("Error") + "] " + txt
	fmt.Println(msg)
}

func Fatal(txt string, output ...string) {
	msg := "[" + color.BoldRed("Fatal error") + "] " + txt
	fmt.Println(msg)
}

func Debug(obj ...interface{}) {
	for i, el := range obj {
		msg := "[" + color.BoldRed("Debug") + "]"
		fmt.Println(msg, i, el)
	}
}

func Bold(txt string) string {
	txt = color.BoldWhite(txt)
	return txt
}

// Prints an event
func Event(event *types.Event) {
	status := "pending"
	if event.Cmd != nil {
		if event.Class == "command_out" {
			if event.Cmd.Status == "error" {
				status = color.BoldRed("Error")
				fmt.Println("  |->", status, event.Cmd.Trace.Errors[0].Error)
				return
			} else if event.Cmd.Status == "success" {
				status = "  |-> " + color.Green("Success")
				var msg string
				var maxWords = 13
				var maxLines = 5
				for i, val := range event.Cmd.ReturnValues {
					if i == maxLines {
						return
					}
					words := strings.Fields(val.(string))
					line := ""
					for ii, v := range words {
						if ii < maxWords {
							line = line + " " + v
						} else {
							line = line + " (...)"
							break
						}
					}
					if i == 0 {
						msg = status + " " + line
					} else {
						msg = msg + line
					}
					if i < len(event.Cmd.ReturnValues)-1 {
						msg = msg + "\n"
					}
					//msg := fmt.Sprintf(" %.120s ", val)
					fmt.Println(msg)
				}
			}
		} else if event.Class == "command_in" {
			msg := " => " + color.Blue("Incoming command") + " " + event.Msg
			fmt.Println(msg)
		} else {
			if event.Cmd.From == "cli" {
				msg := event.Msg
				endMsg := "[" + color.Blue(event.Class) + "] " + msg
				fmt.Println(endMsg)
			}
		}
	} else if event.Class == "state" {
		State(event.Msg)
	} else if event.Class == "status" {
		Status(event.Msg)
	} else if event.Class == "error" {
		Error(event.Msg + "\n" + event.Trace.Msg())
	} else if event.Class == "fatal" {
		Fatal(event.Msg + "\n" + event.Trace.Msg())
	}

}

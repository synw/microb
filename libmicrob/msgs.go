package libmicrob

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
)

func Warning(txt string) {
	msg := "[" + color.Magenta("Warning") + "] " + txt
	if checkDisplay() == true {
		fmt.Println(msg)
	}
}

func Status(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.Blue("Status") + "] " + txt
		fmt.Println(msg)
	}
}

func State(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.Yellow("State") + "] " + txt
		fmt.Println(msg)
	}
}

func Ready(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.Green("Ready") + "] " + txt
		fmt.Println(msg)
	}
}

func Msg(txt string) {
	if checkDisplay() == true {
		fmt.Println(txt)
	}
}

func Ok(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.Green("Ok") + "] " + txt
		fmt.Println(msg)
	}
}

func Timeout(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.BoldRed("Timeout") + "] " + txt
		fmt.Println(msg)
	}
}

func Error(txt string) {
	if checkDisplay() == true {
		msg := "[" + color.BoldRed("Error") + "] " + txt
		fmt.Println(msg)
	}
}

func Debug(obj ...interface{}) {
	for i, el := range obj {
		msg := "[" + color.BoldRed("Debug") + "]"
		fmt.Println(msg, i, el)
	}
}

func PrintEvent(event *types.Event) {
	if checkDisplay() == true {
		if event.Class == "state" {
			State(event.Msg)
		} else if event.Class == "status" {
			Status(event.Msg)
		} else if event.Class == "command_in" {
			msg := " => " + color.Blue("Incoming command") + " " + event.Msg
			fmt.Println(msg)
		} else if event.Class == "command_out" {
			status := event.Cmd.Status
			if status == "error" {
				status = color.BoldRed("error")
				if Verbose() {
					fmt.Println("    |->", status, event.Cmd.Trace.Format())
				}
			} else if status == "success" {
				status = color.Green("success")
				if Verbose() {
					fmt.Println("    |->", status, event.Cmd.ReturnValues)
				}
			}
		} else {
			msg := "[" + color.Blue(event.Class) + "] " + event.Msg
			fmt.Println(msg)
		}
	}
}

func Bold(txt string) string {
	txt = color.BoldWhite(txt)
	return txt
}

func checkDisplay() bool {
	if Verbose() == true {
		return true
	}
	return false
}

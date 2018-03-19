package msgs

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
)

func Warning(txt string) {
	msg := "[" + color.Magenta("Warning") + "] " + txt
	fmt.Println(msg)
}

func Status(txt string) {
	msg := "[" + color.Blue("Status") + "] " + txt
	fmt.Println(msg)
}

func State(txt string) {
	msg := "[" + color.Yellow("State") + "] " + txt
	fmt.Println(msg)
}

func Ready(txt string) {
	msg := "[" + color.BoldGreen("Ready") + "] " + txt
	fmt.Println(msg)
}

func Msg(txt string) {
	fmt.Println(txt)
}

func Ok(txt string) {
	msg := "[" + color.Green("Ok") + "] " + txt
	fmt.Println(msg)
}

func Timeout(txt string) {
	msg := "[" + color.BoldRed("Timeout") + "] " + txt
	fmt.Println(msg)
}

func Error(txt string) {
	msg := "[" + color.BoldRed("Error") + "] " + txt
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

func PrintEvent(event *types.Event) {
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
			status = color.BoldRed("Error")
			fmt.Println("    |->", status, event.Cmd.Trace.Format())
		} else if status == "success" {
			status = color.Green("Success")
			fmt.Println("    |->", status, event.Cmd.ReturnValues)
		}
	} else {
		msg := "[" + color.Blue(event.Class) + "] " + event.Msg
		fmt.Println(msg)
	}

}

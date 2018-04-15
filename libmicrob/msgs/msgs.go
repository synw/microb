package msgs

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
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

func Tr(tr *terr.Trace) {
	tr.Printc()
}

func Debug(obj ...interface{}) {
	for i, el := range obj {
		msg := "[" + color.BoldRed("Debug") + "]"
		fmt.Println(msg, i, el)
	}
}

func Bold(txt string, output ...string) string {
	txt = color.BoldWhite(txt)
	return txt
}

func PrintEvent(event *types.Event) {
	if event.Class == "state" {
		State(event.Msg)
	} else if event.Class == "status" {
		Status(event.Msg)
	} else if event.Class == "error" {
		Error(event.Msg)
		if event.Trace != nil {
			Tr(event.Trace)
		}
	} else if event.Class == "command_in" {
		if event.Cmd.From == "cli" {
			msg := " => " + color.Blue("Incoming command") + " " + event.Msg
			fmt.Println(msg)
		}
	} else if event.Class == "command_out" {
		status := event.Cmd.Status
		if status == "error" {
			if event.Cmd.From == "cli" {
				status = color.BoldRed("Error")
				fmt.Println("  |->", status, event.Cmd.Trace.Format())
			}
		} else if status == "success" {
			if event.Cmd.From == "cli" {
				status = color.Green("Success")
				fmt.Println("  |->", status, event.Cmd.ReturnValues)
			}
		}
	} else {
		if event.Cmd.From == "cli" {
			msg := "[" + color.Blue(event.Class) + "] " + event.Msg
			fmt.Println(msg)
		}
	}

}

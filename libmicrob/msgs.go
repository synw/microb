package libmicrob

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
)

func Warning(txt string, disp ...bool) {
	display := checkDisplay(disp)
	msg := "[" + color.Magenta("Warning") + "] " + txt
	if display == true {
		fmt.Println(msg)
	}
}

func Status(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.Blue("Status") + "] " + txt
		fmt.Println(msg)
	}
}

func State(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.Yellow("State") + "] " + txt
		fmt.Println(msg)
	}
}

func Ready(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.Green("Ready") + "] " + txt
		fmt.Println(msg)
	}
}

func Msg(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		fmt.Println(txt)
	}
}

func Ok(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.Green("Ok") + "] " + txt
		fmt.Println(msg)
	}
}

func Timeout(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.BoldRed("Timeout") + "] " + txt
		fmt.Println(msg)
	}
}

func Error(txt string, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		msg := "[" + color.BoldRed("Error") + "] " + txt
		fmt.Println(msg)
	}
}

func Debug(obj interface{}) {
	msg := "[" + color.BoldRed("Debug") + "]"
	fmt.Println(msg, obj)
}

func PrintEvent(event *types.Event, disp ...bool) {
	display := checkDisplay(disp)
	if display == true {
		if event.Class == "state" {
			State(event.Msg)
		} else if event.Class == "status" {
			Status(event.Msg)
		} else {
			msg := "[" + color.Blue(event.Class) + "]" + event.Msg
			fmt.Println(msg)
		}
	}
}

func Bold(txt string) string {
	txt = color.BoldWhite(txt)
	return txt
}

func checkDisplay(display []bool) bool {
	if len(display) > 0 {
		// To be able to use these functions in the client
		// without relying on the current Microb instance's state
		return true
	} else {
		if Verbose() == true {
			return true
		}
	}
	return false
}

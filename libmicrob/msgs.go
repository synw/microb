package libmicrob

import (
	"fmt"
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb/libmicrob/types"
)

func State(txt string) {
	if Verbose() == true {
		msg := "[" + color.Yellow("State") + "] " + txt
		fmt.Println(msg)
	}
}

func Status(txt string) {
	if Verbose() == true {
		msg := "[" + color.Blue("Status") + "] " + txt
		fmt.Println(msg)
	}
}

func Msg(txt string) {
	if Verbose() == true {
		fmt.Println(txt)
	}
}

func Debug(obj interface{}) {
	msg := "[" + color.Red("Debug") + "]"
	fmt.Println(msg, obj)
}

func Ok(txt string) {
	if Verbose() == true {
		msg := "[" + color.Green("Ok") + "] " + txt
		fmt.Println(msg)
	}
}

func PrintEvent(event *types.Event) {
	if Verbose() == true {
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

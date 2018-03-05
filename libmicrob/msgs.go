package libmicrob

import (
	"fmt"
	color "github.com/acmacalister/skittles"
)

func State(txt string) {
	if Verb.State.Current() != "zero" {
		msg := "[" + color.Yellow("State") + "] " + txt
		fmt.Println(msg)
	}
}

func Debug(obj interface{}) {
	msg := "[" + color.Red("Debug") + "]"
	fmt.Println(msg, obj)
}

func Ok(txt string) {
	if Verb.State.Current() != "zero" {
		msg := "[" + color.Green("Ok") + "] " + txt
		fmt.Println(msg)
	}
}

func Print(txt string, class string) {
	if Verb.State.Current() != "zero" {
		if class == "state" {
			State(txt)
		} else {
			msg := "[" + color.Blue(class) + "]" + txt
			fmt.Println(msg)
		}
	}
}

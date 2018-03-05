package libmicrob

import (
	"fmt"
	color "github.com/acmacalister/skittles"
)

func State(txt string) {
	if Verbose() == true {
		msg := "[" + color.Yellow("State") + "] " + txt
		fmt.Println(msg)
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

func Print(txt string, class ...string) {
	if Verbose() == true {
		if len(class) > 0 {
			if class[0] == "state" {
				State(txt)
			}
		} else {
			msg := "[" + color.Blue(class) + "]" + txt
			fmt.Println(msg)
		}
	}
}

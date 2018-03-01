package msgs

import (
	"fmt"
	color "github.com/acmacalister/skittles"
)

func State(txt string) {
	msg := "[" + color.Yellow("State") + "] " + txt
	fmt.Println(msg)
}

func Debug(obj interface{}) {
	msg := "[" + color.Yellow("Debug") + "]"
	fmt.Println(msg, obj)
}

package methods

import (
	"strconv"
	"github.com/acmacalister/skittles"
)

func FormatStatusCode(sc int) string {
	var sc_str string
	if sc == 404 {
		sc_str = skittles.BoldRed(strconv.Itoa(sc))
	} else if sc == 200 {
		sc_str = skittles.Green(strconv.Itoa(sc))
	} else if sc == 500 {
		sc_str = skittles.BoldRed(strconv.Itoa(sc))
	}
	return sc_str
}

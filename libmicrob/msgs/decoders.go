package msgs

import (
	gc "github.com/PuerkitoBio/goquery"
	color "github.com/acmacalister/skittles"
	//"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/terr"
	"strings"
)

func Decode(txt string, output ...string) (string, *terr.Trace) {
	encoding := "text"
	if len(output) > 0 {
		encoding = output[0]
	}
	var ntxt string
	if encoding == "html" {
		strings.Replace(txt, "<bold>", "<b>", -1)
		strings.Replace(txt, "</bold>", "</b>", -1)
	} else if encoding == "terminal" {
		r := strings.NewReader(txt)
		doc, err := gc.NewDocumentFromReader(r)
		if err != nil {
			tr := terr.New("msgs.decoders.Decode", err.Error())
			return "", tr
		}
		found := false
		doc.Find("bold").Each(func(i int, s *gc.Selection) {
			ntxt = ntxt + color.BoldWhite(s.Text())
			found = true
		})
		if found == false {
			ntxt = txt
		}
	} else {
		return txt, nil
	}
	return ntxt, nil
}

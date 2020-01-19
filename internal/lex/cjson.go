package lex

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
)

var CJSON chroma.Formatter

type ColoredElement struct {
	Content    string `"json":"content"`
	Whitespace bool   `"json":"whitespace"`
	Comment    bool   `"json":"comment"`
	Red        uint8  `"json":"red"`
	Green      uint8  `"json":"green"`
	Blue       uint8  `"json":"blue"`
	Underline  bool   `"json":"underline"`
	Italic     bool   `"json":"italic"`
	Bold       bool   `"json":"bold"`
}

// This is a slightly modified version of github.com/alecthomas/chroma/blob/master/formatters/json.go
func RegisterCJson() {
	CJSON = formatters.Register("cjson", chroma.FormatterFunc(func(w io.Writer, s *chroma.Style, it chroma.Iterator) error {
		fmt.Fprintln(w, "[")
		i := 0
		for t := it(); t != chroma.EOF; t = it() {
			if i > 0 {
				fmt.Fprintln(w, ",")
			}
			i++
			format := s.Get(t.Type)
			colour := format.Colour
			var red uint8
			var green uint8
			var blue uint8
			if colour.IsSet() {
				red = colour.Red()
				green = colour.Green()
				blue = colour.Blue()
			} else {
				//default to black
				red = 0
				green = 0
				blue = 0
			}
			bytes, err := json.Marshal(ColoredElement{
				Content:    t.String(),
				Whitespace: t.Type == chroma.Text,
				Comment:    t.Type.InCategory(chroma.Comment),
				Red:        red,
				Green:      green,
				Blue:       blue,
				Underline:  format.Underline != chroma.Pass,
				Italic:     format.Italic != chroma.Pass,
				Bold:       format.Bold != chroma.Pass,
			})
			if err != nil {
				return err
			}
			if _, err := fmt.Fprint(w, "  "+string(bytes)); err != nil {
				return err
			}
		}
		fmt.Fprintln(w)
		fmt.Fprintln(w, "]")
		return nil
	}))
}

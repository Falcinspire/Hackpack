package lex

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func FromFile(path string, style string) []*ColoredElement {
	lexer := lexers.Match(filepath.Base(path))
	formatter := formatters.Get("cjson")
	theme := styles.Get(style)
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		panic(err)
	}
	var writer bytes.Buffer
	err = formatter.Format(&writer, theme, iterator)
	if err != nil {
		panic(err)
	}
	data := make([]*ColoredElement, 0)
	err = json.Unmarshal(writer.Bytes(), &data)
	if err != nil {
		panic(err)
	}
	data = stealNewlines(data)
	return data
}

func stealNewlines(elements []*ColoredElement) []*ColoredElement {
	nelements := make([]*ColoredElement, 0)
	for _, element := range elements {
		if !element.Whitespace && strings.HasSuffix(element.Content, "\n") {
			element.Content = strings.TrimRight(element.Content, "\r\n")
			nelements = append(nelements, element, &ColoredElement{"\n", true, false, 0, 0, 0, false, false, false})
		} else {
			nelements = append(nelements, element)
		}
	}
	return nelements
}

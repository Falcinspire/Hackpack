package post

import (
	"strings"

	"github.com/falcinspire/hackpackpdf/internal/lex"
)

const ROW_MAX_CHARS = 60

func BreakIntoLines(elements []*lex.ColoredElement) [][]*lex.ColoredElement {
	res := make([][]*lex.ColoredElement, 0)
	cur := make([]*lex.ColoredElement, 0)
	for _, element := range elements {
		if element.Whitespace == true && (strings.Contains(element.Content, "\n") || strings.Contains(element.Content, "\r")) {
			element.Content = cropAtNewlines(element.Content)
			res = append(res, cur)
			cur = make([]*lex.ColoredElement, 0)
			if element.Content != "" {
				cur = append(cur, element)
			}
		} else {
			cur = append(cur, element)
		}
	}
	res = append(res, cur)
	return res
}

func Process(elements [][]*lex.ColoredElement, skipComments bool) [][]*lex.ColoredElement {
	res := make([][]*lex.ColoredElement, 0)
	for _, line := range elements {
		newLine := make([]*lex.ColoredElement, 0)
		for i := range line {
			if skipComments && line[i].Comment {
				continue
			}
			if line[i].Whitespace {
				line[i].Content = shortenSpaces(line[i].Content)
				line[i].Content = capLength(line[i].Content, ROW_MAX_CHARS)
			}
			newLine = append(newLine, line[i])
		}
		allWS := true
		for _, element := range newLine {
			if !element.Whitespace {
				allWS = false
				break
			}
		}
		if allWS {
			continue
		}
		if newLine[0].Content == "import" || newLine[0].Content == "package" || newLine[0].Content == "#" {
			continue
		}
		res = append(res, newLine)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func shortenSpaces(content string) string {
	content = strings.Replace(content, "    ", " ", -1)
	content = strings.Replace(content, "\t", " ", -1)
	return content
}

func capLength(content string, cap int) string {
	return content[0:min(cap, len(content))]
}

func cropAtNewlines(content string) string {
	i := strings.LastIndex(content, "\r")
	i = max(i, strings.LastIndex(content, "\n"))
	if i == len(content)-1 {
		return ""
	} else {
		return content[i+1:]
	}
}

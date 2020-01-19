package reformat

import (
	"strings"
	"unicode"
)

func Name(camelCase string) string {
	var sb strings.Builder
	pastStart := false
	isFirst := true
	for _, char := range camelCase {
		if pastStart {
			if unicode.IsUpper(char) && isFirst {
				sb.WriteRune(' ')
			}
		} else {
			pastStart = true
		}
		if char == '_' {
			sb.WriteRune(' ')
			isFirst = true
		} else {
			if isFirst {
				sb.WriteRune(unicode.ToUpper(char))
				isFirst = false
			} else {
				sb.WriteRune(char)
			}
		}
	}
	return sb.String()
}

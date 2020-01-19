package reformat

import (
	"fmt"
	"testing"
)

func TestReformat(t *testing.T) {
	assert("ReformatCamelCase", "Reformat Camel Case", t)
	assert("Reformat", "Reformat", t)
	assert("reformat", "Reformat", t)
	assert("reformat_case", "Reformat Case", t)
	assert("dijkstra", "Dijkstra", t)
}

func assert(input, expected string, t *testing.T) {
	result := Name(input)
	fmt.Println(input + " " + result)
	if result != expected {
		t.Errorf("\"%s\" -> \"%s\", should be \"%s\"\n", input, result, expected)
	}
}

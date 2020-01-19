package main

import (
	"github.com/falcinspire/hackpack/internal/lex"
	"github.com/falcinspire/hackpack/internal/service"
)

func main() {
	lex.RegisterCJson()

	service.WebService()
}

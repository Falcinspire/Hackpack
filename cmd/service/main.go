package main

import (
	"github.com/falcinspire/hackpackpdf/internal/lex"
	"github.com/falcinspire/hackpackpdf/internal/service"
)

func main() {
	lex.RegisterCJson()

	service.WebService()
}

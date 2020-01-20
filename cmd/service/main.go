package main

import (
	"fmt"

	"github.com/falcinspire/hackpackpdf/internal/lex"
	"github.com/falcinspire/hackpackpdf/internal/service"
)

func main() {
	fmt.Println("Starting webservice...")
	lex.RegisterCJson()

	service.WebService()
}

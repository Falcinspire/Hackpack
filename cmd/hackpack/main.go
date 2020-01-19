package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/falcinspire/hackpack/internal/app"
	"github.com/falcinspire/hackpack/internal/lex"
	"github.com/falcinspire/hackpack/internal/settings"
)

func main() {
	lex.RegisterCJson()

	config := flag.String("config", "hackpack", "Location of the hackpack.toml file")
	output := flag.String("output", "hackpack", "Location of the generated pdf file")
	init := flag.Bool("init", false, "Generate a default configuration file")
	forceInit := flag.Bool("force-init", false, "Generate a default configuration; overriting the existing one if necessary")
	flag.Parse()

	configPath := *config + ".toml"
	outputPath := *output + ".pdf"

	if *forceInit {
		err := settings.WriteDefault(configPath)
		if err != nil {
			panic(err)
		}
	} else if *init {
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println(configPath + " exists! Delete it to generate a fresh one")
		} else {
			err = settings.WriteDefault(configPath)
			if err != nil {
				panic(err)
			}
		}
	} else {
		err := app.CompileHackpack(configPath, outputPath)
		if err != nil {
			panic(err)
		}
	}
}

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/falcinspire/hackpackpdf/internal/build"
	"github.com/falcinspire/hackpackpdf/internal/fileprovide"
	"github.com/falcinspire/hackpackpdf/internal/lex"
	"github.com/falcinspire/hackpackpdf/internal/post"
	"github.com/falcinspire/hackpackpdf/internal/reformat"
	"github.com/falcinspire/hackpackpdf/internal/settings"
)

func CompileHackpack(configPath, outputPath string) error {
	config := settings.GetDefault()
	if _, err := os.Stat(configPath); err == nil {
		err = settings.Read(configPath, config) //TODO Error handling?
		if err != nil {
			return err
		}
	}
	builder := build.New(
		config.Title,
		config.PageLayout.Code,
		config.PageLayout.Header,
		config.PageLayout.Index,
		config.PageLayout.Columns,
		config.PageLayout.Page,
	)

	ignore, err := compileRegexes(config.Source.Ignore)
	if err != nil {
		return err
	}

	sourceset := make([]fileprovide.Source, 0)
	for _, root := range config.Source.Roots {
		fmt.Println("Looking for files in " + root)
		sourceset = append(sourceset, fileprovide.Match(root, ignore)...)
	}
	fmt.Println("Sorting files")
	sort.Slice(sourceset, func(i, j int) bool {
		return strings.Compare(sourceset[i].Relative, sourceset[j].Relative) < 0
	})
	for _, source := range sourceset {
		fmt.Println("Lexing " + source.Relative)
		lexedSource := lex.FromFile(filepath.Join(source.Root, source.Relative), config.Theme)
		fmt.Println("Processing " + source.Relative)
		lexedLines := post.BreakIntoLines(lexedSource)
		processedLines := post.Process(lexedLines, true)
		sourceTitle := reformat.Name(removeExtension(source.Relative))
		fmt.Println("Writing " + source.Relative)
		builder.AppendTitle(sourceTitle, filepath.Join(filepath.Dir(source.Relative), sourceTitle))
		for _, line := range processedLines {
			builder.AppendCode(line)
		}
		builder.AppendLine()
	}
	fmt.Println("Pages written: " + strconv.Itoa(builder.PageNo()) + "; writing index")
	builder.AppendIndex()
	fmt.Println("Saving")
	builder.SaveAndClose(outputPath)

	return nil
}

func compileRegexes(raws []string) ([]*regexp.Regexp, error) {
	ignore := make([]*regexp.Regexp, len(raws))
	for i, reg := range raws {
		var err error
		fmt.Println("Compiling regex " + reg)
		ignore[i], err = regexp.Compile(reg)
		if err != nil {
			return []*regexp.Regexp{}, err
		}
	}
	return ignore, nil
}

func removeExtension(path string) string {
	return strings.TrimRight(filepath.Base(path), filepath.Ext(path))
}

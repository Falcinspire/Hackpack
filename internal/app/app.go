package app

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/falcinspire/hackpack/internal/build"
	"github.com/falcinspire/hackpack/internal/fileprovide"
	"github.com/falcinspire/hackpack/internal/lex"
	"github.com/falcinspire/hackpack/internal/post"
	"github.com/falcinspire/hackpack/internal/reformat"
	"github.com/falcinspire/hackpack/internal/settings"
)

func CompileHackpack(configPath, outputPath string) error {
	var config *settings.Settings
	if _, err := os.Stat(configPath); err == nil {
		config, err = settings.Read(configPath) //TODO Error handling?
		if err != nil {
			return err
		}
	} else {
		config = settings.GetDefaultSettings()
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

	sourceset := fileprovide.Match(config.Source.Roots, ignore)
	for _, source := range sourceset {
		lexedSource := lex.FromFile(filepath.Join(source.Root, source.Relative), config.Theme)
		lexedLines := post.BreakIntoLines(lexedSource)
		processedLines := post.Process(lexedLines, true)
		sourceTitle := reformat.Name(removeExtension(source.Relative))
		builder.AppendTitle(sourceTitle, filepath.Join(filepath.Dir(source.Relative), sourceTitle))
		for _, line := range processedLines {
			builder.AppendCode(line)
		}
		builder.AppendLine()
	}
	builder.AppendIndex()
	builder.SaveAndClose(outputPath)

	return nil
}

func compileRegexes(raws []string) ([]*regexp.Regexp, error) {
	ignore := make([]*regexp.Regexp, len(raws))
	for i, reg := range raws {
		var err error
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

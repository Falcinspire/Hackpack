package fileprovide

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
)

type Source struct {
	Relative string
	Root     string
}

func Match(roots []string, ignoreExt []*regexp.Regexp) []Source {
	paths := make([]Source, 0)
	for _, root := range roots {
		walker := MatchWalker{lexers.Names(true), ignoreExt, make([]string, 0), root}
		filepath.Walk(root, walker.Walk)
		for _, source := range walker.paths {
			paths = append(paths, Source{source, root})
		}
	}
	sort.Slice(paths, func(i, j int) bool {
		return strings.Compare(paths[i].Relative, paths[j].Relative) < 0
	})
	return paths
}

type MatchWalker struct {
	supportedExt []string
	ignoreExt    []*regexp.Regexp
	paths        []string
	root         string
}

func (walker *MatchWalker) supportedExtension(ext string) bool {
	for _, acc := range walker.supportedExt {
		if ext == acc {
			return true
		}
	}
	return false
}

func (walker *MatchWalker) dontIgnore(path string) bool {
	for _, ignore := range walker.ignoreExt {
		if ignore.Match([]byte(path)) {
			return false
		}
	}
	return true
}

func (walker *MatchWalker) Walk(path string, info os.FileInfo, err error) error {
	ext := filepath.Ext(path)
	if ext == "" {
		return nil
	}
	// if regexMatch.MatchString(path) {
	// 	return nil
	// }
	if walker.supportedExtension(ext[1:]) && walker.dontIgnore(path) {
		rel, err := filepath.Rel(walker.root, path)
		if err != nil {
			panic(err)
		}
		walker.paths = append(walker.paths, rel)
	}
	return nil
}

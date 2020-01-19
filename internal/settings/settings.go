package settings

import (
	"io/ioutil"
	"os"

	"github.com/falcinspire/hackpackpdf/internal/resource"
	"github.com/pelletier/go-toml"
)

type Settings struct {
	Title      string     `toml:"title"`
	Theme      string     `toml:"theme"`
	Source     Source     `toml:"source"`
	PageLayout PageLayout `toml:"page_layout"`
}

type Source struct {
	Roots  []string `toml:"roots"`
	Ignore []string `toml:"ignore"`
}

type PageLayout struct {
	Page    string  `toml:"page"`
	Columns int     `toml:"columns"`
	Code    float64 `toml:"code_font_size"`
	Header  float64 `toml:"header_font_size"`
	Index   float64 `toml:"index_font_size"`
}

func Read(path string) (*Settings, error) {
	settings := GetDefaultSettings()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(bytes, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func GetDefaultSettings() *Settings {
	return &Settings{
		Title: "Hackpack",
		Theme: "colorful",
		Source: Source{
			Roots:  []string{"."},
			Ignore: []string{".*?\\.xml", ".*?\\.json", ".*?\\.toml"},
		},
		PageLayout: PageLayout{
			Page:    "Letter",
			Columns: 3,
			Code:    6.0,
			Header:  10.0,
			Index:   8.0,
		},
	}
}

func WriteDefault(path string) error {
	return ioutil.WriteFile(path, resource.ReadDefaultHackpackToml(), 0644) //TODO magic number
}

func Write(path string, config *Settings) error {
	bytes, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644) //TODO magic number
}

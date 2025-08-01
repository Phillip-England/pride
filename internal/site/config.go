package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Config struct {
	SiteName      string
	Dir           string
	Path          string
	RelativePath  string
	Text          string
	Dob           string
	Version       string
	Server        string
	Theme         string
	ContentDir    string
	TemplatesDir  string
	NavigationDir string
}

func NewConfig(path string) Config {
	var config Config
	config.Path = path
	config.Dir = strings.TrimSuffix(path, "/pride.toml")
	config.SiteName = filepath.Base(config.Dir)
	config.Dob = time.Now().UTC().Format(time.RFC3339)
	config.Version = "0.0.1"
	config.Server = "https://www.example.com"
	config.Theme = "dracula"
	config.Text = fmt.Sprintf(`version = "%s"
dob = "%s"
server = "%s"`, config.Version, config.Dob, config.Server)
	return config
}

func LoadConfig() (Config, *syserr.Err) {
	var config Config
	foundConfig := false
	dir, err := os.Getwd()
	if err != nil {
		return config, syserr.New(syserr.Here(), err.Error())
	}
	for {
		configPath := dir + "/pride.toml"
		file, err := os.ReadFile(configPath)
		if err != nil {
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
			continue
		}
		config.Text = string(file)
		config.Path = configPath
		if _, err := toml.DecodeFile(config.Path, &config); err != nil {
			return config, syserr.New(syserr.Here(), err.Error())
		}
		foundConfig = true
		config.Dir = dir
		config.SiteName = filepath.Base(config.Dir)
		break
	}
	if !foundConfig {
		return config, syserr.New(syserr.Here(), err.Error())
	}
	config.ContentDir = config.Dir + "/content"
	config.TemplatesDir = config.Dir + "/templates"
	config.NavigationDir = config.Dir + "/navigation"
	return config, nil
}

func (config Config) Create() *syserr.Err {
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.Here(), err.Error())
	}
	defer file.Close()
	_, err = file.Write([]byte(config.Text))
	if err != nil {
		return syserr.New(syserr.Here(), err.Error())
	}
	return nil
}

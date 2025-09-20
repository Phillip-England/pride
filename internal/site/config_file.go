package site

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Phillip-England/pride/internal/syserr"
)

type ConfigFile struct {
	Path         string
	RelativePath string
	Text         string
	Dob          string
	Version      string
	Server       string
	Theme        string
}

func CreateConfigFile(path string) (ConfigFile, error) {
	var config ConfigFile
	config.Path = path
	config.Dob = time.Now().UTC().Format(time.RFC3339)
	config.Version = "0.0.1"
	config.Server = "https://www.example.com"
	config.Theme = "dracula"
	config.Text = fmt.Sprintf(`version = "%s"
dob = "%s"
server = "%s"
theme = "%s"`, config.Version, config.Dob, config.Server, config.Theme)
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return config, syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()

	_, err = file.Write([]byte(config.Text))
	if err != nil {
		return config, syserr.New(syserr.Here(), "%s", err.Error())
	}

	return config, nil
}

func LoadConfigFile() (ConfigFile, error) {
	var config ConfigFile
	foundConfig := false
	dir, err := os.Getwd()
	if err != nil {
		return config, syserr.New(syserr.Here(), "%s", err.Error())
	}
	checkedPaths := []string{}
	for {
		configPath := filepath.Join(dir, "pride.toml")
		checkedPaths = append(checkedPaths, configPath)
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
			return config, syserr.New(syserr.Here(), "%s", err.Error())
		}
		foundConfig = true
		break
	}
	if !foundConfig {
		return config, syserr.New(syserr.Here(), "failed to locate pride.toml at any of the following locations:\n %v", checkedPaths)
	}
	return config, nil
}

func (config ConfigFile) Create() error {
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer file.Close()
	_, err = file.Write([]byte(config.Text))
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	return nil
}

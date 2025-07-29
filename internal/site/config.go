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

func ConfigNew(path string) Config {
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

func ConfigLoadFromCwd() (Config, syserr.SysErr) {
	var config Config
	configPath := "./pride.toml"
	_, err := os.Stat(configPath)
	if err != nil {
		return config, syserr.New(syserr.CodeMia, fmt.Errorf("failed to located %s", configPath))
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, syserr.New(syserr.CodeMia, fmt.Errorf("failed to read %s", configPath))
	}
	config.Text = string(file)
	config.Path = configPath
	if _, err := toml.DecodeFile(config.Path, &config); err != nil {
		return config, syserr.New(syserr.CodeDev, fmt.Errorf("failed to decode %s, at this point, the file has been verified as existing, so this is a developer error", config.Path))
	}
	return config, nil
}

func ConfigLoad() (Config, syserr.SysErr) {
	var config Config
	foundConfig := false
	dir, err := os.Getwd()
	if err != nil {
		return config, syserr.DevNew(fmt.Errorf("failed to get a reference to the cwd becuase:\n%s", err.Error()))
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
			return config, syserr.New(syserr.CodeDev, fmt.Errorf("failed to decode %s, at this point, the file has been verified as existing, so this is a developer error", config.Path))
		}
		foundConfig = true
		config.Dir = dir
		config.SiteName = filepath.Base(config.Dir)
		break
	}
	if !foundConfig {
		return config, syserr.MiaNew(fmt.Errorf("failed to load pride.toml from current dir, or any of it's parent dirs"))
	}
	config.ContentDir = config.Dir + "/content"
	config.TemplatesDir = config.Dir + "/templates"
	config.NavigationDir = config.Dir + "/navigation"
	return config, nil
}

func (config Config) Create() syserr.SysErr {
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when creating %s", config.Path))
	}
	defer file.Close()
	_, err = file.Write([]byte(config.Text))
	if err != nil {
		return syserr.New(syserr.CodeHelp, fmt.Errorf("unanticipated error when writing to %s", config.Path))
	}
	return nil
}

package site

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Config struct {
	Path    string
	Text    string
	Dob     string
	Version string
	Server  string
}

func ConfigNew(path string) Config {
	var config Config
	config.Path = path
	config.Dob = time.Now().UTC().Format(time.RFC3339)
	config.Version = "0.0.1"
	config.Server = "https://www.example.com"
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

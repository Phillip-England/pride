package pride

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type FileConfig struct {
	Path    string
	Text    string
	Dob     string
	Version string
	Server  string
}

func FileConfigNew(path string) FileConfig {
	var config FileConfig
	config.Path = path
	config.Dob = time.Now().UTC().Format(time.RFC3339)
	config.Version = "0.0.1"
	config.Server = "https://www.example.com"
	config.Text = fmt.Sprintf(`version = "%s"
dob = "%s"
server = "%s"`, config.Version, config.Dob, config.Server)
	return config
}

func FileConfigLoadFromCwd() (FileConfig, SysErr) {
	var config FileConfig
	configPath := "./pride.toml"
	_, err := os.Stat(configPath)
	if err != nil {
		return config, SysErrNew(SysErrCodeMia, fmt.Errorf("failed to located %s", configPath))
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, SysErrNew(SysErrCodeMia, fmt.Errorf("failed to read %s", configPath))
	}
	config.Text = string(file)
	config.Path = configPath
	if _, err := toml.DecodeFile(config.Path, &config); err != nil {
		return config, SysErrNew(SysErrCodeDev, fmt.Errorf("failed to decode %s, at this point, the file has been verified as existing, so this is a developer error", config.Path))
	}
	return config, nil
}

func (config FileConfig) Create() SysErr {
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when creating %s", config.Path))
	}
	defer file.Close()
	_, err = file.Write([]byte(config.Text))
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when writing to %s", config.Path))
	}
	return nil
}

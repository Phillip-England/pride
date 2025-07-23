package pride

import (
	"fmt"
	"os"
	"time"
)

type FileConfig struct {
	Path string
	Text string
	Dob  string
}

func FileConfigNew(path string) FileConfig {
	var config FileConfig
	config.Path = path
	config.Dob = time.Now().UTC().Format(time.RFC3339)
	config.Text = fmt.Sprintf(`version = "0.0.1"
dob = "%s"
server = "https://www.example.com"`, config.Dob)
	return config
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

package pride

import (
	"fmt"
	"os"
)

type FileTemplateBase struct {
	Path string
	Text string
}

func FileTemplateBaseNew(path string) FileTemplateBase {
	var config FileTemplateBase
	config.Path = path
	config.Text = `{{ define "header" }}
	<header>
		<h1>{{ .Pride.Ctx.Title }}</h1>
	</header>
{{ end }}
	
{{ define "nav" }}
	{{ .Pride.Nav.main }}
{{ end }}

{{ define "footer" }}
	<footer>
		<p>Created: {{ .Pride.Ctx.dob }}</p>
	</footer>
{{ end }}

`
	return config
}

func (f FileTemplateBase) Create() SysErr {
	file, err := os.OpenFile(f.Path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when creating %s", f.Path))
	}
	defer file.Close()
	_, err = file.Write([]byte(f.Text))
	if err != nil {
		return SysErrNew(SysErrCodeHelp, fmt.Errorf("unanticipated error when writing to %s", f.Path))
	}
	return nil
}

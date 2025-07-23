package pride

import (
	"fmt"
	"os"
)

type FileLayoutBase struct {
	Path string
	Text string
}

func FileLayoutBaseNew(path string) FileLayoutBase {
	var config FileLayoutBase
	config.Path = path
	config.Text = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="stylesheet" href="/static/default_styles.css">
		<title>{{ .Pride.Ctx.Title }}</title>
	</head>
  	<body>
		<div id='root'>
			{{ template "header" . }}
			{{ template "nav" . }}
			<main>{{ .Pride.Content }}</main>
			{{ template "footer" . }}
		</div>
  	</body>
</html>`
	return config
}

func (f FileLayoutBase) Create() SysErr {
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

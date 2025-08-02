package syserr

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

type Err struct {
	File    string
	Line    int
	Message string
	Err     error
}

func New(location *Location, format string, args ...any) *Err {
	var serr Err
	message := fmt.Sprintf(format, args...)
	serr.File = location.File
	serr.Line = location.Line
	serr.Message = message
	serr.Err = errors.New(message)
	return &serr
}

func (e Err) Print() {
	fmt.Fprintf(os.Stderr, "🚨 %s:%d — %s\n", e.File, e.Line, e.Message)
}

func (e Err) Fail() {
	fmt.Fprintf(os.Stderr, "🚨 %s:%d — %s\n", e.File, e.Line, e.Message)
	os.Exit(1)
}

func (e Err) Error() string {
	return fmt.Sprintf("perr.Err [%s:%d] %s", e.File, e.Line, e.Message)
}

type Location struct {
	File string
	Line int
}

func Here() *Location {
	_, file, line, _ := runtime.Caller(1)
	return &Location{File: file, Line: line}
}

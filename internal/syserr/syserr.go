package syserr

import (
	"errors"
	"fmt"
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
	fmt.Printf("🚨 perr.Err\n")
	fmt.Printf("  File   : %s\n", e.File)
	fmt.Printf("  Line   : %d\n", e.Line)
	fmt.Printf("  Message: %s\n", e.Message)
	if e.Err != nil {
		fmt.Printf("  Cause  : %v\n", e.Err)
	}
}

func (e Err) Error() string {
	return fmt.Sprintf("perr.Err [%s:%d] code=%d: %s", e.File, e.Line, e.Message)
}

type Location struct {
	File string
	Line int
}

func Here() *Location {
	_, file, line, _ := runtime.Caller(1)
	return &Location{File: file, Line: line}
}

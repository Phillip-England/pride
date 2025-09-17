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
    message := fmt.Sprintf(format, args...)
    return &Err{
        File:    location.File,
        Line:    location.Line,
        Message: message,
        Err:     errors.New(message),
    }
}

func (e *Err) Error() string {
    return fmt.Sprintf("syserr [%s:%d] %s", e.File, e.Line, e.Message)
}

func (e *Err) Unwrap() error {
    return e.Err
}

func (e *Err) Print() {
    fmt.Fprintf(os.Stderr, "🚨 %s:%d — %s\n", e.File, e.Line, e.Message)
}

func (e *Err) Fail() {
    e.Print()
    os.Exit(1)
}

type Location struct {
    File string
    Line int
}

func Here() *Location {
    _, file, line, _ := runtime.Caller(1)
    return &Location{File: file, Line: line}
}

func Consume(location *Location, err error) *Err {
	if err == nil {
		return nil
	}
	if serr, ok := err.(*Err); ok {
		return &Err{
			File:    location.File,
			Line:    location.Line,
			Message: serr.Message,
			Err:     serr,
		}
	}
	return &Err{
		File:    location.File,
		Line:    location.Line,
		Message: err.Error(),
		Err:     err,
	}
}

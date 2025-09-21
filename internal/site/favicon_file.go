package site

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/Phillip-England/pride/internal/syserr"
)

//go:embed favicon.ico
var favicon []byte

type FaviconFile struct {
	Path  string
	Bytes []byte
}

// CreateFaviconFile writes the embedded favicon bytes to the given path and returns the FaviconFile struct.
func CreateFaviconFile(path string) (FaviconFile, error) {
	var file FaviconFile
	file.Bytes = favicon
	file.Path = path
	if err := os.WriteFile(path, favicon, 0644); err != nil {
		return file, syserr.New(syserr.Here(), "%s", fmt.Errorf("failed to write favicon file: %w", err))
	}
	return file, nil
}

// LoadFaviconFile checks if the favicon file exists and loads its metadata.
func LoadFaviconFile(path string) (FaviconFile, error) {
	var file FaviconFile
	file.Path = path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return file, syserr.New(syserr.Here(), "%s", fmt.Errorf("favicon file does not exist at path: %s", path))
	} else if err != nil {
		return file, syserr.New(syserr.Here(), "%s", fmt.Errorf("error checking favicon file: %w", err))
	}
	return file, nil
}

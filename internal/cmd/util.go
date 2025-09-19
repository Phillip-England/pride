// Package cmd provides utility functions for parsing and handling
// command-line arguments in a Go application.
package cmd

import (
	"fmt"
	"os"
	"slices"
)

// GetArg returns the nth command-line argument from os.Args.
// If the requested index is out of range, it returns an error.
//
// Parameters:
//   n - the index of the argument to retrieve
//
// Returns:
//   string - the argument value if it exists
//   error  - non-nil if n is out of bounds
func GetArg(n int) (string, error) {
	if len(os.Args) <= n {
		return "", fmt.Errorf(`out of bounds access on os.Args`)
	}
	return os.Args[n], nil
}

// ArgIsFilePath checks whether the argument at the given position can be
// treated as a valid file path. It attempts to create the file with write
// permissions and immediately removes it to verify access.
//
// Parameters:
//   position - the index of the argument to check
//
// Returns:
//   string - the argument value if valid
//   bool   - true if the argument is a writable file path, false otherwise
func ArgIsFilePath(position int) (string, bool) {
	arg, err := GetArg(position)
	if err != nil {
		return "", false
	}
	f, err := os.OpenFile(arg, os.O_CREATE, 0644)
	if err != nil {
		return "", false
	}
	f.Close()
	os.Remove(arg)
	return arg, true
}

// HasFlag checks whether a specific flag exists in the command-line arguments.
// Uses slices.Contains for a clean, modern implementation.
//
// Parameters:
//   flag - the command-line flag to search for
//
// Returns:
//   bool - true if the flag exists, false otherwise
func HasFlag(flag string) bool {
    return slices.Contains(os.Args, flag)
}

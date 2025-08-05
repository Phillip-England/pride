package cmd

import (
	"fmt"
	"os"
)

func GetArg(n int) (string, error) {
	if len(os.Args) <= n {
		return "", fmt.Errorf(`out of bounds access on os.Args`)
	}
	return os.Args[n], nil
}

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

func HasFlag(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

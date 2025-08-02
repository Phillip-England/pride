package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Phillip-England/pride/internal/site"
)

func testDir() string {
	path, _ := filepath.Abs("./tmp/test-site")
	return path
}

func clean() {
	_ = os.RemoveAll(testDir())
}

func TestMain(t *testing.T) {
	clean()
	dir, serr := site.NewPrideDir(testDir())
	if serr != nil {
		serr.Fail()
		return
	}
	serr = dir.CreateIfNotExists()
	if serr != nil {
		serr.Fail()
	}

}
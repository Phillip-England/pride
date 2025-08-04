package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Phillip-England/pride/internal/site"
)

func testSitePath() string {
	wd, _ := os.Getwd()
	path, _ := filepath.Abs(filepath.Join(wd, "tmp", "test-site"))
	return path
}

func clean() {
	_ = os.RemoveAll(testSitePath())
}

func TestPrideDirCreationAndLoading(t *testing.T) {
	clean()
	testSitePath := testSitePath()
	_, serr := site.NewPrideDir(testSitePath)
	if serr != nil {
		serr.Fail()
		return
	}
	deepDir := filepath.Join(testSitePath, "content", "posts", "foo", "bar")
	err := os.MkdirAll(deepDir, 0755)
	if err != nil {
		t.Fatal(err.Error())
	}	
	err = os.Chdir(deepDir)
	if err != nil {
		t.Fatal(err.Error())
	}
	prideDir, serr := site.LoadPrideDir()
	if serr != nil {
		serr.Fail()
		return
	}
	if prideDir.Path != testSitePath {
		t.Fatalf("the pride directory has an unexpected rootDir\nhave: %s\nexpect: %s", prideDir.Path, testSitePath)
	}
}
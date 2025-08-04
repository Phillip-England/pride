package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Phillip-England/pride/internal/op"
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

func initTestSite() site.PrideDir {
	clean()
	dir, _ := site.CreatePrideDir(testSitePath())
	_ = os.Chdir(dir.Path)
	return dir
}

// ensures the pride dir can be created and accessed
// tests to see if the pride dir can be loaded from any subdirectory
// within the pride dir
func TestCreatePrideDir(t *testing.T) {
	clean()
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	testSitePath := testSitePath()
	_, serr := site.CreatePrideDir(testSitePath)
	if serr != nil {
		serr.Fail()
		return
	}
	deepDir := filepath.Join(testSitePath, "content", "posts", "foo", "bar")
	err = os.MkdirAll(deepDir, 0755)
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
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestOperationNewContent(t *testing.T) {
	dir := initTestSite()
	contentPath := filepath.Join(dir.ContentDir.Path, "about.md")
	serr := op.OperationNewContent(contentPath)
	if serr != nil {
		serr.Fail()
		return
	}

}

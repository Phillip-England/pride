package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

// ensures new content can be generated
func TestOperationNewContent(t *testing.T) {
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	dir := initTestSite()
	contentPath := filepath.Join(dir.ContentDir.Path, "about.md")
	mdFile, serr := op.OperationNewContent(contentPath, true)
	if serr != nil {
		serr.Fail()
		return
	}
	if mdFile.ServerPath != "/about" {
		t.Fatalf("expected server path /about but got %s", mdFile.ServerPath)
	}
	if mdFile.Title != "About" {
		t.Fatalf("expected title About but got %s", mdFile.Title)
	}
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBuildNavigation(t *testing.T) {
	dir := initTestSite()
	contentPaths := []string{
		filepath.Join(dir.ContentDir.Path, "about.md"),
		filepath.Join(dir.ContentDir.Path, "contact.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/1.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/2.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/3.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/1.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/2.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/3.md"),
	}
	for _, path := range contentPaths {
		_, serr := op.OperationNewContent(path, true)
		if serr != nil {
			serr.Fail()
			return
		}
	}
	time.Sleep(1)
	op.BuildNavigation()
}

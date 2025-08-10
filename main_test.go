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

func initBlankTestSite() site.PrideDir {
	clean()
	dir, _ := site.CreatePrideDir(testSitePath())
	_ = os.Chdir(dir.Path)
	return dir
}

func initFullTestSite() site.PrideDir {
	clean()
	dir, _ := site.CreatePrideDir(testSitePath())
	_ = os.Chdir(dir.Path)
	contentPaths := []string{
		filepath.Join(dir.ContentDir.Path, "contact.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/a-good-post.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/a-funny-post.md"),
		filepath.Join(dir.ContentDir.Path, "/posts/that-post.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/these-docs.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/those-docs.md"),
		filepath.Join(dir.ContentDir.Path, "/docs/them-docs.md"),
	}
	for _, path := range contentPaths {
		_, _ = op.OperationNewContent(path, true)
	}
	prideDir, _ := site.LoadPrideDir()
	return prideDir
}

// 1. ensures pride dir can be generated
// 2. ensures pride dir cannot be overwritten
// 3. ensures pride commands can be executed deep within the pride dir
// 4. ensures commands cannot be executed outside of a pride dir
func TestCreatePrideDir(t *testing.T) {
	clean()
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	testSitePath := testSitePath()
	// 1
	_, serr := site.CreatePrideDir(testSitePath)
	if serr != nil {
		serr.Fail()
		return
	}
	// 2
	_, serr = site.CreatePrideDir(testSitePath)
	if serr == nil {
		t.Fatal("pride dir was able to be overwritten, which is invalid behaviour")
		return
	}
	// 3
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
	// 4
	prideDir, serr = site.LoadPrideDir()
	if serr == nil {
		t.Fatal("was able to execute a pride command outside of a pride dir, which is invalid behaviour")
		return
	}
}

// 1. ensures new content can be generated
// 2. ensures new content is generated with the expected server path
// 3. ensures new content is generated with the expected title
// 4. ensures content with complex, hyphenated names can be generated
// 5. ensures content with complex, hyphenated names is generated with the expected server path
// 6. ensures content with complex, hyphenated names is generated with the expected title
// 7. ensures existing content cannot be overwritten
// 8. ensures we can create content in bulk without any glaring issues
func TestOperationNewContent(t *testing.T) {
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	dir := initBlankTestSite()
	contentPath := filepath.Join(dir.ContentDir.Path, "about.md")
	// 1
	mdFile, serr := op.OperationNewContent(contentPath, true)
	if serr != nil {
		serr.Fail()
		return
	}
	// 2
	if mdFile.ServerPath != "/about" {
		t.Fatalf("expected server path '/about' but got '%s'", mdFile.ServerPath)
	}
	// 3
	if mdFile.Title != "About" {
		t.Fatalf("expected title 'About' but got '%s'", mdFile.Title)
	}
	// 4
	complexContentPath := filepath.Join(dir.ContentDir.Path, "some-difficult-name-that-might-resolve-weird.md")
	complexMdFile, serr := op.OperationNewContent(complexContentPath, true)
	if serr != nil {
		serr.Fail()
		return
	}
	// 5
	if complexMdFile.ServerPath != "/some-difficult-name-that-might-resolve-weird" {
		t.Fatalf("expected server path '/some-difficult-name-that-might-resolve-weird' but got '%s'", complexMdFile.ServerPath)
	}
	// 6
	if complexMdFile.Title != "Some Difficult Name that Might Resolve Weird" {
		t.Fatalf("expected title 'Some Difficult Name that Might Resolve Weird' but got '%s'", complexMdFile.Title)
	}
	// 7
	mdFile, serr = op.OperationNewContent(contentPath, true)
	if serr == nil {
		t.Fatal("was able to overwrite content which already exists, which is invalid behaviour")
		return
	}
	// 8
	contentPaths := []string{
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
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// 1. content is loaded in initFullTestSite()
// 2. ensuring we have 8 .md files
func TestLoadContent(t *testing.T) {
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	// 1
	prideDir := initFullTestSite()
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
	// 2
	mdFiles := prideDir.ContentDir.MarkdownFiles
	if len(mdFiles) != 8 {
		t.Fatalf(`expected 8 .md files but found %d`, len(mdFiles))
	}
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
}

// 1. navigation is loaded in initFullTestSite()
func TestLoadNavigation(t *testing.T) {
	startingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}
	// 1
	_ = initFullTestSite()
	err = os.Chdir(startingDir)
	if err != nil {
		t.Fatal(err.Error())
	}
}

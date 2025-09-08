package build

import (
	"os"

	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Build struct {
	Src site.PrideDir
	Dest string
}

func dirExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    return false
}


func GenerateBuild(dest string) (*syserr.Err) {
	// if the dest already exists, exit
	err := dirExists(dest)
	if err {
		return syserr.New(syserr.Here(), "<DESTINATION> %s already exists", dest)
	}
	// load up the prideDir
	prideDir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	// load up the server
	svr, serr := server.NewServer(8080, prideDir)
	if serr != nil {
		return serr
	}
	// go through each route, generating an html file for each one
	for _, route := range svr.Routes {
		_, serr := NewHtmlFile(dest, route)
		if serr != nil {
			return serr
		}
		// fmt.Println(dest, htmlFile)
	}
	return nil
}


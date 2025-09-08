package build

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
)

type Build struct {
	Src site.PrideDir
}

func GenerateBuild() (*syserr.Err) {
	fmt.Println("building...")
	return nil
}


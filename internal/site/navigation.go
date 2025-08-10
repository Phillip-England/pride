package site

import (
	"fmt"

	"github.com/Phillip-England/pride/internal/syserr"
)

type Navigation struct {
	Menus map[string]Menu
}

func LoadNavigation(contentDir ContentDir) (Navigation, *syserr.Err) {
	var nav Navigation
	fmt.Println("hit")
	return nav, nil
}

type Menu struct {
	Name  string
	Links []Link
}

type Link struct {
	Href string
	Html string
}

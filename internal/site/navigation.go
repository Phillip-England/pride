package site

import (
	"github.com/Phillip-England/pride/internal/syserr"
)

type Navigation struct {
	Menus map[string]Menu
}

func LoadNavigation(contentDir ContentDir) (Navigation, *syserr.Err) {
	var nav Navigation
	menus, serr := LoadMenus(contentDir)
	if serr != nil {
		return nav, serr
	}
	nav.Menus = menus
	return nav, nil
}

type Menu struct {
	Name  string
	Path  string
	Links []Link
	Html  string
}

func LoadMenus(contentDir ContentDir) (map[string]Menu, *syserr.Err) {
	menus := make(map[string]Menu)
	return menus, nil
}

type Link struct {
	Href string
	Html string
}

func LoadLinks(path string, contentDir ContentDir) ([]Link, *syserr.Err) {
	links := []Link{}
	return links, nil
}

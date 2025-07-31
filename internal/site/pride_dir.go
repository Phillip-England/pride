package site

import "path/filepath"

type PrideDir struct {
	Root          string
	ContentDir    string
	ConfigPath    string
	StaticDir     string
	TemplatesDir  string
	NavigationDir string
	LayoutsDir    string
}

// Returns a fully populated PrideDir struct based on a root path
func NewPrideDir(path string) *PrideDir {
	return &PrideDir{
		Root:          path,
		ContentDir:    filepath.Join(path, "content"),
		ConfigPath:    filepath.Join(path, "pride.toml"),
		StaticDir:     filepath.Join(path, "static"),
		TemplatesDir:  filepath.Join(path, "templates"),
		NavigationDir: filepath.Join(path, "navigation"),
		LayoutsDir:    filepath.Join(path, "layouts"),
	}
}

package main

import (
	"embed"
	"io/fs"
	"path/filepath"
	"text/template"
)

//go:embed "pages"
var Pages embed.FS

//go:embed "static"
var Static embed.FS

func loadTemplates() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(Pages, "pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		ts, err := template.ParseFS(Pages, page)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		cache[name] = ts
	}
	return cache, nil
}

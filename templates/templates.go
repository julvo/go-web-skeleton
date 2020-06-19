package templates

import (
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var isDev = false

func init() {
	isDev = strings.HasPrefix(strings.ToUpper(os.Getenv("ENV")), "DEV")
}

type Templates struct {
	dir       string
	templates map[string]*template.Template
}

func New(dir string) (Templates, error) {
	tmplMap, err := LoadTemplateMap(dir)
	if err != nil {
		return Templates{}, err
	}

	return Templates{
		dir:       dir,
		templates: tmplMap,
	}, nil
}

func (t *Templates) Reload() error {
	tmplMap, err := LoadTemplateMap(t.dir)
	if err != nil {
		return err
	}

	t.templates = tmplMap
	return nil
}

func (t *Templates) Execute(templName string, w io.Writer, data interface{}) error {
	if isDev {
		if err := t.Reload(); err != nil {
			return err
		}
	}

	tmpl, ok := t.templates[templName]
	if !ok {
		return errors.New("Could not find template named " + templName)
	}

	return tmpl.Execute(w, data)
}

func LoadTemplateMap(dir string) (map[string]*template.Template, error) {
	partials, err := filepath.Glob(filepath.Join(dir, "partials", "*.html"))
	if err != nil {
		return nil, err
	}
	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.html"))
	if err != nil {
		return nil, err
	}
	root := filepath.Join(dir, "root.html")

	templates := make(map[string]*template.Template)
	for _, page := range pages {
		relPath, _ := filepath.Rel(filepath.Join(dir, "pages"), page)
		templ, err := template.ParseFiles(append([]string{root, page}, partials...)...)
		if err != nil {
			return nil, err
		}
		templates[relPath] = templ
	}
	return templates, nil
}

package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

type Templates struct {
	templates *template.Template
}

// Render lookup & execute template to HTML document, then minify that document
func (t *Templates) Render(w io.Writer, name string, data interface{}, _ echo.Context) (err error) {
	wTmpl := &bytes.Buffer{}
	if err = t.templates.ExecuteTemplate(wTmpl, name, data); err != nil {
		return
	}
	if err = minifier.Minify(echo.MIMETextHTML, w, bytes.NewReader(wTmpl.Bytes())); err != nil {
		return
	}
	return
}

func scanTemplateFiles(root string) ([]string, error) {
	var filepaths []string

	err := filepath.Walk(root, func(path string, finfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !finfo.IsDir() && filepath.Ext(finfo.Name()) == ".gohtml" {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	return filepaths, err
}

func newTemplates() *Templates {
	filepaths, err := scanTemplateFiles("./views")
	if err != nil {
		panic(err)
	}
	return &Templates{
		templates: template.Must(template.ParseFiles(filepaths...)),
	}
}

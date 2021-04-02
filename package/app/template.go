package app

import (
	"go/build"
	"html/template"
)

func GetTemplate() *template.Template {
	path := build.Default.GOPATH + "/src/CRUD/templates/*"

	tmpl := template.Must(template.ParseGlob(path))
	return tmpl
}

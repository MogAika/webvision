package views

import "html/template"

var Templates *template.Template

func InitTemplates() (err error) {
	Templates, err = template.ParseFiles(
		"templates/index.html",
		"templates/errors.html",
	)
	return
}

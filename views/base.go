package views

import (
	"html/template"
	"net/http"
)

var Templates *template.Template

func InitTemplates() (err error) {
	Templates, err = template.ParseFiles(
		"templates/index.html",
		"templates/errors.html",
		"templates/browse.html",
		"templates/upload.html",
	)
	return
}

func ViewError(w http.ResponseWriter, code int, title, text string) {
	type ErrorView struct {
		Code  int
		Title string
		Text  template.HTML
	}

	w.WriteHeader(code)

	Templates.ExecuteTemplate(w, "error", &ErrorView{
		Code:  code,
		Title: title,
		Text:  template.HTML(text),
	})
}

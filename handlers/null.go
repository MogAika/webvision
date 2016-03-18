package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/mogaika/webvision/views"
)

func HandlerNotFound(w http.ResponseWriter, r *http.Request) {
	ViewError(w, 404, "Not found", r.URL.String())
}

func HandlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf(`URL: %v<br>
	Form: %#v<br>
	MultipartForm: %#v<br>
	PostForm: %#v<br>
	RemoteAddr: %v`, r.URL, r.Form, r.MultipartForm, r.PostForm, r.RemoteAddr)

	ViewError(w, 500, "Not implemented", s)
}

func ViewError(w http.ResponseWriter, code int, title, text string) {
	type ErrorView struct {
		Code  int
		Title string
		Text  template.HTML
	}

	views.Templates.ExecuteTemplate(w, "error", &ErrorView{
		Code:  code,
		Title: title,
		Text:  template.HTML(text),
	})
}

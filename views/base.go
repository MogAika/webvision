package views

import (
	"html/template"
	"net/http"
)

var Templates *template.Template

func InitTemplates() (err error) {
	Templates, err = template.ParseFiles(
		"templates/globals.html",
		"templates/errors.html",
		"templates/browse.html",
		"templates/upload.html",
		"templates/random.html",
	)
	return
}

type ViewErrorStruct struct {
	Code  int
	Title string
	Text  template.HTML
}

func ViewError(w http.ResponseWriter, code int, title, text string) {
	w.WriteHeader(code)

	Templates.ExecuteTemplate(w, "error", &ViewErrorStruct{
		Code:  code,
		Title: title,
		Text:  template.HTML(text),
	})
}

func ViewBrowse(w http.ResponseWriter) {
	Templates.ExecuteTemplate(w, "browse", nil)
}

func ViewRandom(w http.ResponseWriter) {
	Templates.ExecuteTemplate(w, "random", nil)
}

type ViewUploadStruct struct {
	StorageUsageMessage string
	StorageUsagePercent float32
}

func ViewUpload(w http.ResponseWriter, datausage string, datausagepercent float32) {
	Templates.ExecuteTemplate(w, "upload", &ViewUploadStruct{
		StorageUsageMessage: datausage,
		StorageUsagePercent: datausagepercent,
	})
}

func ViewUploadResult(w http.ResponseWriter, err string) {
	w.Write([]byte(err))
}

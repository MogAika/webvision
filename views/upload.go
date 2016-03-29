package views

import "net/http"

func ViewUpload(w http.ResponseWriter) {
	Templates.ExecuteTemplate(w, "upload", nil)
}

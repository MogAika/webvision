package views

import "net/http"

func ViewUpload(w http.ResponseWriter) {
	Templates.ExecuteTemplate(w, "upload", nil)
}

func ViewUploadResult(w http.ResponseWriter, err string) {
	w.Write([]byte(err))
}

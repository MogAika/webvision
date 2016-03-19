package views

import "net/http"

func ViewBrowse(w http.ResponseWriter) {
	Templates.ExecuteTemplate(w, "browse", nil)
}

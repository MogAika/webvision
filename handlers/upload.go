package handlers

import (
	"fmt"
	"net/http"

	"github.com/mogaika/webvision/views"
)

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf(`Form: %#v<br>
	MultipartForm: %#v<br>
	PostForm: %#v<br>`, r.Form, r.MultipartForm, r.PostForm)

	f, fh, err := r.FormFile("heh")

	s = fmt.Sprintf(`Heh file: %#v<br>
			Heh head: %#v<br>
			Heh error: %#v<br>`, f, fh, err) + s

	if err == nil {
		defer f.Close()
		views.ViewError(w, 200, "Uploaded", s)
	} else {
		views.ViewError(w, 500, "Not implemented", s)
	}
}

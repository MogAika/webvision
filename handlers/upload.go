package handlers

import (
	"net/http"

	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

func HandlerUploadGet(w http.ResponseWriter, r *http.Request) {
	views.ViewUpload(w)
}

func HandlerUploadPost(w http.ResponseWriter, r *http.Request) {
	db, set := VarsFromRequest(r)
	r.Body = http.MaxBytesReader(w, r.Body, set.MaxDataSize)

	f, fh, err := r.FormFile("fl")

	if err == nil {
		defer f.Close()
		_, err := (&models.Media{}).NewFromFile(db, f, fh.Header.Get("Content-Type"), set)
		if err == nil {
			views.ViewUploadResult(w, "")
		} else {
			views.ViewUploadResult(w, err.Error())
		}
	} else {
		views.ViewError(w, 500, "Error", err.Error())
	}
}

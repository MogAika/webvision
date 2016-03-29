package handlers

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/views"
)

func ProcessFile(f multipart.File) error {
	buffer := make([]byte, 512)

	n, readerr := f.Read(buffer)
	if readerr != nil && readerr != io.EOF {
		return readerr
	}

	ctype := http.DetectContentType(buffer[:n])
	log.Log.Info(ctype)

	return nil
}

func HandlerUploadGet(w http.ResponseWriter, r *http.Request) {
	views.ViewUpload(w)
}

func HandlerUploadPost(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("heh")

	if err == nil {
		defer f.Close()
		err = ProcessFile(f)
		if err == nil {
			views.ViewError(w, 200, "Uploaded", fh.Filename+" uploaded")
		} else {
			log.Log.Errorf("Error processing file \"%s\": %v", fh.Filename, err)
			views.ViewError(w, 500, "Error", "Error processing file on server side")
		}
	} else {
		views.ViewError(w, 500, "Error", err.Error())
	}
}

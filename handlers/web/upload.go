package web

import (
	"net/http"

	"github.com/mogaika/webvision/views"
)

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	views.ViewUpload(w)
}

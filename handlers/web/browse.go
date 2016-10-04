package web

import (
	"net/http"

	"github.com/mogaika/webvision/views"
)

func HandlerBrowse(w http.ResponseWriter, r *http.Request) {
	views.ViewBrowse(w)
}

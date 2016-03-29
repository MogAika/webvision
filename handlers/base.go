package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/settings"
	"github.com/mogaika/webvision/views"
)

func VarsFromRequest(r *http.Request) (*gorm.DB, *settings.Settings) {
	return context.Get(r, "db").(*gorm.DB), context.Get(r, "settings").(*settings.Settings)
}

func HandlerNotFound(w http.ResponseWriter, r *http.Request) {
	views.ViewError(w, 404, "Not found", r.URL.String())
}

func HandlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf(`URL: %v<br>
	Form: %#v<br>
	MultipartForm: %#v<br>
	PostForm: %#v<br>
	RemoteAddr: %v`, r.URL, r.Form, r.MultipartForm, r.PostForm, r.RemoteAddr)

	views.ViewError(w, 500, "Not implemented", s)
}

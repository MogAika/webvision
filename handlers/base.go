package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/views"
)

func VarsFromRequest(r *http.Request) (*gorm.DB, *sessions.CookieStore) {
	return context.Get(r, "db").(*gorm.DB), context.Get(r, "cookiestore").(*sessions.CookieStore)
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

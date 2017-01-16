package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mogaika/webvision/handlers/api"
	"github.com/mogaika/webvision/handlers/web"
	"github.com/mogaika/webvision/helpers"
)

func AuthedUserHandler(redirect bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, conf := helpers.ContextGetVars(r.Context())
		if conf.Secret != "" && !helpers.UserIsAuthorized(r) {
			if redirect {
				http.SetCookie(w, &http.Cookie{
					Name:  "redirect",
					Value: r.RequestURI,
					Path:  "/login",
				})
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	}
}

func InitRouter(r *mux.Router) {
	r.HandleFunc("/", AuthedUserHandler(true, web.HandlerBrowse))
	r.HandleFunc("/login", web.HandlerLogin)
	r.HandleFunc("/upload", AuthedUserHandler(true, web.HandlerUpload))
	r.HandleFunc("/random", AuthedUserHandler(true, web.HandlerRandom))

	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/login", api.HandlerApiLogin)
	a.HandleFunc("/query", AuthedUserHandler(false, api.HandlerApiQuery))
	a.HandleFunc("/upload", AuthedUserHandler(false, api.HandlerApiUpload))
	a.HandleFunc("/random", AuthedUserHandler(false, api.HandlerApiRandom))
}

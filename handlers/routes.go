package handlers

import (
	"github.com/gorilla/mux"

	"github.com/mogaika/webvision/handlers/api"
	"github.com/mogaika/webvision/handlers/web"
)

func InitRouter(r *mux.Router) {
	r.HandleFunc("/", web.HandlerBrowse)
	r.HandleFunc("/upload", web.HandlerUpload)
	r.HandleFunc("/random", web.HandlerRandom)

	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/query", api.HandlerApiQuery)
	a.HandleFunc("/random", api.HandlerApiRandom)
	a.HandleFunc("/upload", api.HandlerApiUpload)
}

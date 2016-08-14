package handlers

import (
	"net/http"
	"path"

	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

func HandlerRandom(w http.ResponseWriter, r *http.Request) {
	db, set := VarsFromRequest(r)

	md, err := (&models.Media{}).GetRandom(db)
	if err != nil {
		views.ViewError(w, 500, "Internal error", err.Error())
	} else {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		http.ServeFile(w, r, path.Join(set.DataPath, md.File))
	}
}

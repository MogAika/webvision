package web

import (
	"fmt"
	"math/rand"
	"net/http"
	"path"

	"github.com/ricochet2200/go-disk-usage/du"

	"github.com/mogaika/webvision/helpers"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

func HandlerRandom(w http.ResponseWriter, r *http.Request) {
	db, conf := helpers.ContextGetVars(r.Context())

	md, err := (&models.Media{}).GetRandom(db, rand.Int31())
	if err != nil {
		views.ViewError(w, 500, "Internal error", err.Error())
	} else {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		http.ServeFile(w, r, path.Join(conf.DataPath, md.File))
	}
}

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	_, conf := helpers.ContextGetVars(r.Context())

	usage := du.NewDiskUsage(conf.DataPath)

	fmt.Fprintf(w, `Storage use %vM/%vM (%v%%)`,
		usage.Used()/1024/1024, usage.Size()/1024/1024, usage.Usage()*100.0)
}

func HandlerBrowse(w http.ResponseWriter, r *http.Request) {
	views.ViewBrowse(w)
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	views.ViewUpload(w)
}

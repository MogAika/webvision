package web

import (
	"fmt"
	"net/http"

	"github.com/ricochet2200/go-disk-usage/du"

	"github.com/mogaika/webvision/helpers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/views"
)

func HandlerRandom(w http.ResponseWriter, r *http.Request) {
	views.ViewRandom(w)
}

func HandlerBrowse(w http.ResponseWriter, r *http.Request) {
	views.ViewBrowse(w)
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Log.Errorf("Error parsing form %v: %v", r.RemoteAddr, err)
		} else {
			_, conf := helpers.ContextGetVars(r.Context())
			if r.FormValue("secret") == conf.Secret {
				if err := helpers.DoAuth(w, r, conf); err != nil {
					w.WriteHeader(503)
					log.Log.Errorf("Error when auth: %v", err)
				} else {
					if redirect, err := r.Cookie("redirect"); err == nil && redirect != nil {
						http.Redirect(w, r, redirect.Value, http.StatusMovedPermanently)
					} else {
						http.Redirect(w, r, "/", http.StatusMovedPermanently)
					}
				}
			}
		}
	}
	views.ViewLogin(w)
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	_, conf := helpers.ContextGetVars(r.Context())

	usage := du.NewDiskUsage(conf.DataPath)

	usagestr := fmt.Sprintf(`Storage usage: %d Mb / %d Mb (%.2f%%)`, usage.Used()/1024/1024, usage.Size()/1024/1024, usage.Usage()*100.0)

	views.ViewUpload(w, usagestr, usage.Usage()*100.0)
}

package web

import (
	"fmt"
	"net/http"

	"github.com/ricochet2200/go-disk-usage/du"

	"github.com/mogaika/webvision/helpers"
	"github.com/mogaika/webvision/views"
)

func HandlerRandom(w http.ResponseWriter, r *http.Request) {
	views.ViewRandom(w)
}

func HandlerBrowse(w http.ResponseWriter, r *http.Request) {
	views.ViewBrowse(w)
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	_, conf := helpers.ContextGetVars(r.Context())

	usage := du.NewDiskUsage(conf.DataPath)

	usagestr := fmt.Sprintf(`Storage usage: %d Mb / %d Mb (%.2f%%)`, usage.Used()/1024/1024, usage.Size()/1024/1024, usage.Usage()*100.0)

	views.ViewUpload(w, usagestr, usage.Usage()*100.0)
}

package web

import (
	"fmt"
	"net/http"

	"github.com/ricochet2200/go-disk-usage/du"

	"github.com/mogaika/webvision/helpers"
)

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	_, conf := helpers.ContextGetVars(r.Context())

	usage := du.NewDiskUsage(conf.DataPath)

	fmt.Fprintf(w, `Storage use %vM/%vM (%v%%)`,
		usage.Used()/1024/1024, usage.Size()/1024/1024, usage.Usage()*100.0)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/ricochet2200/go-disk-usage/du"
)

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	_, set := VarsFromRequest(r)

	usage := du.NewDiskUsage(set.DataPath)

	fmt.Fprintf(w, `Storage use %vM/%vM (%v%%)`,
		usage.Used()/1024/1024, usage.Size()/1024/1024, usage.Usage()*100.0)
}

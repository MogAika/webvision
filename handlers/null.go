package handlers

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	spew.Fprintf(w, "not implemented\n%#v\n", r)
}

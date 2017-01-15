package api

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/mogaika/webvision/helpers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
)

const (
	JSON_CONTENT_TYPE = "application/json; charset=utf-8"
)

type ViewMedia struct {
	Id    uint64
	Url   string
	Thumb *string
	Type  string
}

func (vm *ViewMedia) FromModel(md *models.Media) *ViewMedia {
	vm.Id = md.ID
	vm.Url = md.File
	vm.Thumb = md.Thumbnail
	vm.Type = md.Type
	return vm
}

func HandlerApiQuery(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		log.Log.Criticalf("Error parsing requested url: %v", err)
		return
	}

	q := u.Query()

	start := -1
	count := 25

	if !stoi(q, "start", &start) || !stoi(q, "count", &count) {
		return
	}

	var medias_data []ViewMedia
	var md []models.Media

	db, _ := helpers.ContextGetVars(r.Context())
	if start <= 0 {
		md, err = (&models.Media{}).Get(db, count)
	} else if start == 1 {
		err = nil
	} else {
		md, err = (&models.Media{}).GetTo(db, start, count)
	}

	if err != nil {
		log.Log.Errorf("Error quering db: %v", err)
		return
	}

	medias_data = make([]ViewMedia, len(md))
	for i, v := range md {
		medias_data[i].FromModel(&v)
	}

	apiWrite(w, medias_data)
}

func HandlerApiRandom(w http.ResponseWriter, r *http.Request) {
	db, _ := helpers.ContextGetVars(r.Context())

	md, err := (&models.Media{}).GetRandom(db, rand.Int31())
	if err != nil {
		apiError(w, err.Error())
	} else {
		apiWrite(w, new(ViewMedia).FromModel(md))
	}
}

func HandlerApiUpload(w http.ResponseWriter, r *http.Request) {
	db, conf := helpers.ContextGetVars(r.Context())

	r.Body = http.MaxBytesReader(w, r.Body, conf.MaxDataSize)

	f, fh, err := r.FormFile("fl")

	var ff io.ReadCloser
	ct := ""

	if err != nil {
		urlToUpload := r.FormValue("url")
		if urlToUpload == "" {
			apiError(w, err.Error())
			return
		}

		fh, err := http.Get(urlToUpload)
		if err != nil {
			apiError(w, err.Error())
			return
		}
		ff = fh.Body
		ct = fh.Header.Get("Content-Type")
	} else {
		ff = f
		ct = fh.Header.Get("Content-Type")
	}

	defer ff.Close()
	md, err := (&models.Media{}).NewFromFile(db, ff, ct, conf)
	if err != nil {
		apiError(w, err.Error())
	} else {
		apiWrite(w, new(ViewMedia).FromModel(md))
	}
}

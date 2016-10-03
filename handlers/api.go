package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
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

func stoi(vals url.Values, key string, pInt *int) bool {
	v := vals.Get(key)
	if v != "" {
		ival, err := strconv.Atoi(v)
		if err != nil {
			log.Log.Errorf("Error parsing query val 's': %v", key, err)
			return false
		} else {
			*pInt = ival
		}
	}
	return true
}

func apiWrite(w http.ResponseWriter, data interface{}) {
	binData, err := json.Marshal(data)
	if err != nil {
		log.Log.Errorf("Error marshal responce: %v", err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", JSON_CONTENT_TYPE)

	w.Write(binData)
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

	db, _ := VarsFromRequest(r)
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
	db, _ := VarsFromRequest(r)

	md, err := (&models.Media{}).GetRandom(db, rand.Int31())
	if err != nil {
		views.ViewError(w, 500, "Internal error", err.Error())
	} else {
		apiWrite(w, new(ViewMedia).FromModel(md))
	}
}

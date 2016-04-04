package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

const (
	MEDIAS_PER_REQUEST = 25
)

func HandlerBrowse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-requested-with") == "XMLHttpRequest" {
		// if ajax
		u, err := url.Parse(r.RequestURI)
		if err != nil {
			log.Log.Criticalf("Error parsing requested url: %v", err)
			return
		}

		sstart := u.Query().Get("s")
		start := -1
		if sstart != "" {
			var err error
			start, err = strconv.Atoi(sstart)
			if err != nil {
				log.Log.Errorf("Error parsing query val 's': %v", err)
				return
			}
		}

		type ViewMedia struct {
			Id    uint64
			Url   string
			Thumb *string
			Type  string
		}

		var medias_data []ViewMedia
		var md []models.Media

		db, _ := VarsFromRequest(r)
		if start < 0 {
			md, err = (&models.Media{}).Get(db, MEDIAS_PER_REQUEST)
		} else if start <= 1 {
			err = nil
		} else {
			md, err = (&models.Media{}).GetTo(db, start, MEDIAS_PER_REQUEST)
		}

		if err != nil {
			log.Log.Errorf("Error quering db: %v", err)
			return
		}

		medias_data = make([]ViewMedia, len(md))
		for i, v := range md {
			medias_data[i].Id = v.ID
			medias_data[i].Url = v.File
			medias_data[i].Thumb = v.Thumbnail
			medias_data[i].Type = v.Type
		}

		dat, err := json.Marshal(medias_data)
		if err != nil {
			log.Log.Errorf("Error marshal responce: %v", err)
			return
		}
		w.Write(dat)
	} else {
		views.ViewBrowse(w)
	}
}

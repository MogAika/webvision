package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mogaika/webvision/log"
)

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
		log.Log.Errorf("Error marshal response: %v", err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", JSON_CONTENT_TYPE)

	w.Write(binData)
}

func apiError(w http.ResponseWriter, inerr interface{}) {
	binData, err := json.Marshal(map[string]interface{}{"error": inerr})
	if err != nil {
		log.Log.Errorf("Error marshal response: %v", err)
		return
	}
	w.Write(binData)
}

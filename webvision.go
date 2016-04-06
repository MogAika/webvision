package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mogaika/webvision/app"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/settings"
)

func main() {
	conffile := os.Getenv("WEBVISION_CONFIG")

	if conffile == "" {
		conffile = "webvision.yaml"
	}

	log.Log.Infof("Used config %s", conffile)

	confdata, err := ioutil.ReadFile(conffile)
	if err != nil {
		log.Log.Criticalf("Cannot open config file '%s': %v", conffile, err)
	} else {
		s := &settings.Settings{}
		err = yaml.Unmarshal(confdata, s)

		if err != nil {
			log.Log.Criticalf("Error cannot read yaml config file '%s': %v", conffile, err)
			os.Exit(1)
		} else {
			_, err = app.NewApp(s)
			if err != nil {
				log.Log.Errorf("Application error: %v\n", err)
			}
		}
	}
}

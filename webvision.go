package main

import (
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/mogaika/webvision/app"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/settings"
)

var Config *string

func init() {
	Config = flag.String("c", "webvision.yaml", "Config file")
}

func main() {
	flag.Parse()

	confdata, err := ioutil.ReadFile(*Config)
	if err != nil {
		log.Log.Criticalf("Cannot open config file '%s': %v", *Config, err)
	} else {
		s := &settings.Settings{}
		err = yaml.Unmarshal(confdata, s)

		if err != nil {
			log.Log.Criticalf("Error cannot read yaml config file '%s': %v", *Config, err)
			os.Exit(1)
		} else {
			_, err = app.NewApp(s)
			if err != nil {
				log.Log.Errorf("Application error: %v\n", err)
			}
		}
	}
}

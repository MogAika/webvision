package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/mogaika/webvision/app"
	"github.com/mogaika/webvision/config"
	"github.com/mogaika/webvision/log"
)

func loadConfig(filename string) *config.Config {
	log.Log.Infof("Used config %s", filename)

	confdata, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Log.Criticalf("Cannot open config file '%s': %v", filename, err)
		return nil
	}
	conf := &config.Config{}
	err = yaml.Unmarshal(confdata, conf)

	if err != nil {
		log.Log.Criticalf("Error cannot read yaml config file '%s': %v", filename, err)
		return nil
	}
	return conf

}

func main() {
	var configFileName string

	flag.StringVar(&configFileName, "config", "config.yaml", "Path to configuration file")
	flag.Parse()

	conf := loadConfig(configFileName)
	if conf == nil {
		return
	}

	a, err := app.NewApp(conf)
	if err != nil {
		log.Log.Criticalf("Error creating app: %v", err)
		return
	}

	if err := a.Listen(); err != nil {
		log.Log.Errorf("Server fault: %v", err)
	}
}

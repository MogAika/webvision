package main

import (
	"github.com/mogaika/webvision/app"
	"github.com/mogaika/webvision/log"
)

func main() {
	s := &app.AppSettings{
		DataPath: "./data/",
		Secret:   "aegasnjp9r8hO2da",
		DB: app.AppDBSettings{
			Dialect: "sqlite3",
			Params:  "_test.db",
		},
		Web: app.AppWebSettings{
			Host:        "127.0.0.1:8080",
			Url:         "",
			Tls:         true,
			TlsCertFile: "server.pem",
			TlsKeyFile:  "server.key",
		},
	}

	_, err := app.NewApp(s)
	if err != nil {
		log.Log.Errorf("Application error: %v\n", err)
	}
}

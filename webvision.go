package main

import (
	"github.com/alexcesaro/log/stdlog"

	"github.com/mogaika/webvision/app"
)

func main() {
	log := stdlog.GetFromFlags()

	s := &app.AppSettings{
		DataPath: "./data/",
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

	_, err := app.NewApp(s, log)
	if err != nil {
		log.Errorf("Application error: %v\n", err)
	}
}

package main

import (
	"github.com/mogaika/webvision/app"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/settings"
)

func main() {
	s := &settings.Settings{
		DataPath:    "./data/",
		MaxDataSize: 1024 * 1024 * 128, // 128 Mbytes
		DB: settings.DBSettings{
			Dialect: "sqlite3",
			Params:  "_test.db",
		},
		Web: settings.WebSettings{
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

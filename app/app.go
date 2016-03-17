package app

import (
	_log "log"
	"net/http"
	"os"

	"github.com/alexcesaro/log"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/handlers"
	"github.com/mogaika/webvision/models"
)

type App struct {
	DB  *gorm.DB
	Log log.Logger

	Settings *AppSettings
}

type AppWebSettings struct {
	Host        string // 0.0.0.0:8080
	Url         string // https://www.url.com
	Tls         bool
	TlsCertFile string
	TlsKeyFile  string
}

type AppDBSettings struct {
	Dialect string
	Params  interface{}
}

type AppSettings struct {
	DataPath string
	DB       AppDBSettings
	Web      AppWebSettings
}

func NewApp(s *AppSettings, log log.Logger) (a *App, err error) {
	a = &App{
		Log:      log,
		Settings: s,
	}

	a.DB, err = gorm.Open(s.DB.Dialect, s.DB.Params)
	if err != nil {
		return
	}

	a.DB.SetLogger(_log.New(os.Stdout, "", _log.LstdFlags))
	if log.LogInfo() {
		a.DB.LogMode(true)
	}

	log.Info("Initializing db")
	models.Init(a.DB)

	log.Info("Starting server")

	return a, a.InitHttp()
}

func (a *App) InitHttp() error {
	r := mux.NewRouter()

	if a.Settings.Web.Url != "" {
		r.Host(a.Settings.Web.Url)
	}

	r.HandleFunc("/", handlers.NotImplemented)
	r.HandleFunc("/data/{id:[0-9]+}", handlers.NotImplemented)
	r.HandleFunc("/upload", handlers.NotImplemented)
	r.HandleFunc("/statews", handlers.NotImplemented)
	r.HandleFunc("/statews", handlers.NotImplemented)

	h := gorillahandlers.RecoveryHandler()(r)

	host := a.Settings.Web.Host

	if a.Settings.Web.Tls {
		if host == "" {
			host = ":443"
		}
		http.ListenAndServeTLS(host, a.Settings.Web.TlsCertFile, a.Settings.Web.TlsKeyFile, h)
	} else {
		if host == "" {
			host = ":80"
		}
		http.ListenAndServe(host, h)
	}

	return nil
}

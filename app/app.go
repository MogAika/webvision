package app

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/handlers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/settings"
	"github.com/mogaika/webvision/views"
)

type App struct {
	DB       *gorm.DB
	Settings *settings.Settings
	Handlers http.Handler
}

func NewApp(s *settings.Settings) (a *App, err error) {
	a = &App{
		Settings: s,
	}

	if err = views.InitTemplates(); err != nil {
		return
	}

	a.DB, err = gorm.Open(s.DB.Dialect, s.DB.Params)
	if err != nil {
		return
	}

	a.DB.SetLogger(log.NewGormLogger(log.Log))
	if log.Log.LogInfo() {
		a.DB.LogMode(true)
	}

	log.Log.Info("Initializing db")
	models.Init(a.DB)

	log.Log.Info("Starting server")

	return a, a.InitHttp()
}

func (a *App) Handler(h http.Handler) http.Handler {
	a.Handlers = h
	return a
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.Set(r, "db", a.DB)
	context.Set(r, "settings", a.Settings)
	context.Set(r, "app", a)

	a.Handlers.ServeHTTP(w, r)

	context.Clear(r)
}

func (a *App) InitHttp() error {
	r := mux.NewRouter()

	if a.Settings.Web.Url != "" {
		r.Host(a.Settings.Web.Url)
	}

	r.NotFoundHandler = &NotFoundHandler{}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir(a.Settings.DataPath))))

	r.HandleFunc("/", handlers.HandlerBrowse)
	r.HandleFunc("/upload", handlers.HandlerUploadGet).Methods("GET")
	r.HandleFunc("/upload", handlers.HandlerUploadPost).Methods("POST")

	h := a.Handler(r)

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

type NotFoundHandler struct{}

func (NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerNotFound(w, r)
}

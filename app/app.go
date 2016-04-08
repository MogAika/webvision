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

	r.NotFoundHandler = &NotFoundHandler{}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir(a.Settings.DataPath))))

	r.HandleFunc("/", handlers.HandlerBrowse)
	r.HandleFunc("/upload", handlers.HandlerUploadGet).Methods("GET")
	r.HandleFunc("/upload", handlers.HandlerUploadPost).Methods("POST")
	r.HandleFunc("/status", handlers.HandlerStatus)

	h := a.Handler(r)

	httperror := make(chan error)

	if a.Settings.Web.Host != "" {
		go func() {
			httperror <- http.ListenAndServe(a.Settings.Web.Host, h)
		}()
	}

	if a.Settings.Web.TlsHost != "" {
		go func() {
			httperror <- http.ListenAndServeTLS(a.Settings.Web.TlsHost, a.Settings.Web.TlsCertFile, a.Settings.Web.TlsKeyFile, h)
		}()
	}
	return <-httperror
}

type NotFoundHandler struct{}

func (NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerNotFound(w, r)
}

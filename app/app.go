package app

import (
	"context"
	"net/http"

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
	ctx := context.WithValue(r.Context(), "settings", a.Settings)
	ctx = context.WithValue(ctx, "db", a.DB)
	a.Handlers.ServeHTTP(w, r.WithContext(ctx))
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
	r.HandleFunc("/random", handlers.HandlerRandom)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/query", handlers.HandlerApiQuery)
	api.HandleFunc("/random", handlers.HandlerApiRandom)

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

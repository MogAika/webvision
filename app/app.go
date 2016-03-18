package app

import (
	"net/http"

	"github.com/gorilla/context"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/handlers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

type App struct {
	DB       *gorm.DB
	Settings *AppSettings

	CookieStore *sessions.CookieStore
	Handlers    http.Handler
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
	Secret   string
}

func NewApp(s *AppSettings) (a *App, err error) {
	a = &App{
		Settings:    s,
		CookieStore: sessions.NewCookieStore([]byte(s.Secret)),
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
	context.Set(r, "cookiestore", a.CookieStore)
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

	r.HandleFunc("/", handlers.HandlerNotImplemented)
	r.HandleFunc("/data/{id:[0-9]+}", handlers.HandlerNotImplemented)
	r.HandleFunc("/upload", handlers.HandlerNotImplemented)

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

type NotFoundHandler struct{}

func (NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerNotFound(w, r)
}

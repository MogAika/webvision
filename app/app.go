package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/config"
	"github.com/mogaika/webvision/handlers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

type App struct {
	DB       *gorm.DB
	Config   *config.Config
	Handlers http.Handler
}

func NewApp(conf *config.Config) (a *App, err error) {
	a = &App{Config: conf}

	if err = views.InitTemplates(); err != nil {
		return
	}

	a.DB, err = gorm.Open(conf.DB.Dialect, conf.DB.Params)
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

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "conf", a.Config)
	ctx = context.WithValue(ctx, "db", a.DB)
	a.Handlers.ServeHTTP(w, r.WithContext(ctx))
}

func (a *App) InitHttp() (err error) {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir(a.Config.DataPath))))

	handlers.InitRouter(r)
	a.Handlers = r

	return err
}

func (a *App) Listen() (err error) {
	if a.Config.Web.CertFile != "" && a.Config.Web.KeyFile != "" {
		err = http.ListenAndServeTLS(a.Config.Web.Host, a.Config.Web.CertFile, a.Config.Web.KeyFile, a)
	} else {
		err = http.ListenAndServe(a.Config.Web.Host, a)
	}
	return
}

package app

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/config"
	"github.com/mogaika/webvision/handlers"
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/views"
)

type App struct {
	DB           *gorm.DB
	Config       *config.Config
	Handlers     http.Handler
	SecureCookie *securecookie.SecureCookie
}

func NewApp(conf *config.Config) (a *App, err error) {
	a = &App{Config: conf}

	if err = views.InitTemplates(); err != nil {
		return
	}

	a.SecureCookie = securecookie.New([]byte(conf.Cookie.HashKey), []byte(conf.Cookie.BlockKey))

	log.Log.Info("Initializing db")
	a.DB, err = gorm.Open(conf.DB.Dialect, conf.DB.Params)
	if err != nil {
		return
	}

	a.DB.SetLogger(log.NewGormLogger(log.Log))
	if log.Log.LogInfo() {
		a.DB.LogMode(true)
	}
	models.Init(a.DB)

	log.Log.Info("Checking config")
	a.PrintWarnings()
	log.Log.Info("Starting server")
	return a, a.InitHttp()
}

func (a *App) PrintWarnings() {
	if a.Config.Secret == "" {
		log.Log.Warningf("Site secret is not used")
	}
	hashKeyLen := len(a.Config.Cookie.HashKey)
	if hashKeyLen != 32 && hashKeyLen != 64 {
		log.Log.Warningf("Bad cookie hash key length (%d). Prefer use 32 or 64", hashKeyLen)
	}
	blockKeyLen := len(a.Config.Cookie.BlockKey)
	if blockKeyLen != 16 && blockKeyLen != 24 && blockKeyLen != 32 {
		log.Log.Warningf("Bad cookie block key length (%d). Prefer use 16, 24 or 32", blockKeyLen)
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "conf", a.Config)
	ctx = context.WithValue(ctx, "db", a.DB)
	ctx = context.WithValue(ctx, "cs", a.SecureCookie)

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

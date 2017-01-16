package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/config"
	"github.com/mogaika/webvision/log"
)

func ContextGetVars(ctx context.Context) (*gorm.DB, *config.Config) {
	return ctx.Value("db").(*gorm.DB), ctx.Value("conf").(*config.Config)
}

func ContextGetSecureCookie(ctx context.Context) *securecookie.SecureCookie {
	return ctx.Value("cs").(*securecookie.SecureCookie)
}

func UserIsAuthorized(r *http.Request) bool {
	if authCookie, err := r.Cookie("auth"); err == nil {
		var value string
		if err := ContextGetSecureCookie(r.Context()).Decode("auth", authCookie.Value, &value); err == nil {
			authEstimateTime := &time.Time{}
			if err := authEstimateTime.UnmarshalText([]byte(value)); err == nil {
				return authEstimateTime.After(time.Now())
			} else {
				log.Log.Noticef("Probably cookie chiper attack %v: %v", r.RemoteAddr, err)
			}
		} else {
			log.Log.Noticef("Probably cookie hash attack %v: %v", r.RemoteAddr, err)
		}
	}
	return false
}

func DoAuth(w http.ResponseWriter, r *http.Request, conf *config.Config) error {
	expiresTime := time.Now().Add(time.Duration(conf.Cookie.LifeTime) * time.Second)
	if timeText, err := expiresTime.MarshalText(); err == nil {
		if timeEncoded, err := ContextGetSecureCookie(r.Context()).Encode("auth", string(timeText)); err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:    "auth",
				Value:   timeEncoded,
				Path:    "/",
				Expires: expiresTime,
			})
			return nil
		} else {
			return err
		}
	} else {
		return err
	}

}

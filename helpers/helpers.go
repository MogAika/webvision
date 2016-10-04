package helpers

import (
	"golang.org/x/net/context"

	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/config"
)

func ContextGetVars(ctx context.Context) (*gorm.DB, *config.Config) {
	return ctx.Value("db").(*gorm.DB), ctx.Value("conf").(*config.Config)
}

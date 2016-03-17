package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Media struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Type string  `xorm:"not null varchar(64)"`
	Hash string  `xorm:"not null varchar(24)"`
	Size int64   `xorm:"not null"`
	Path *string `xorm:"varchar(256)"`

	Tags []Tag `gorm:"many2many:m2m_media_tag;"`
}

func (md *Media) New(db *gorm.DB, ftype, fhash string, fsize int64) (*Media, error) {
	md = &Media{
		Type: ftype,
		Hash: fhash,
		Size: fsize,
	}

	return md, db.Create(md).Error
}

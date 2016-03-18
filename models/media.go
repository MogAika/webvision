package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Media struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Type      string  `xorm:"not null varchar(64)"`
	Hash      string  `xorm:"not null varchar(24)"`
	Size      int64   `xorm:"not null"`
	File      *string `xorm:"varchar(256)"`
	Thumbnail *string `xorm:"varchar(256)"`

	Tags []Tag `gorm:"many2many:m2m_media_tag;"`
}

func (md *Media) New(db *gorm.DB, ftype string, fsize int64) (*Media, error) {
	md = &Media{
		Type: ftype,
		Size: fsize,
	}

	return md, db.Create(md).Error
}

func (md *Media) SetFile(db *gorm.DB, file, hash string, thumbnail *string) error {
	upd := map[string]interface{}{
		"file":      &file,
		"hash":      hash,
		"thumbnail": thumbnail,
	}

	return db.Model(md).Updates(upd).Error
}

func (md *Media) AddTag(db *gorm.DB, tag *Tag) error {
	return db.Model(md).Association("Tags").Append(tag).Error
}
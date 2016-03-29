package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Media struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Type      string  `xorm:"not null varchar(64)"`
	Hash      string  `xorm:"not null varchar(24)"`
	Size      int64   `xorm:"not null"`
	File      string  `xorm:"varchar(256)"`
	Thumbnail *string `xorm:"varchar(256)"`
	Likes     int64   `xorm:"not null"`
	Dislikes  int64   `xorm:"not null"`

	Tags []Tag `gorm:"many2many:m2m_media_tag;"`
}

func (md *Media) New(db *gorm.DB, file, hash, ftype string, fsize int64, thumbnail *string) (*Media, error) {
	md = &Media{
		File:      file,
		Hash:      hash,
		Type:      ftype,
		Size:      fsize,
		Thumbnail: thumbnail,
	}
	return md, db.Create(md).Error
}

func (md *Media) TagsAdd(db *gorm.DB, tag *Tag) error {
	return db.Model(md).Association("Tags").Append(tag).Error
}

func (md *Media) TagsGet(db *gorm.DB) (error, []Tag) {
	var tags []Tag
	err := db.Model(md).Association("Tags").Find(&tags).Error
	return err, tags
}

func (md *Media) TagsRemove(db *gorm.DB, tag *Tag) error {
	return db.Model(md).Association("Tags").Delete(tag).Error
}

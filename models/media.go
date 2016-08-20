package models

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/settings"
)

var ErrAlreadyUploaded = errors.New("File already uploaded on server")
var ErrIncorrectContentType = errors.New("Server support only image/video/audio files")

type Media struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time

	Type      string  `xorm:"not null varchar(128)"`
	Hash      string  `xorm:"not null varchar(32)"`
	Size      int64   `xorm:"not null"`
	File      string  `xorm:"varchar(256)"`
	Thumbnail *string `xorm:"varchar(256)"`
	Likes     int64   `xorm:"not null"`
	Dislikes  int64   `xorm:"not null"`
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

func (md *Media) Get(db *gorm.DB, limit int) ([]Media, error) {
	var media []Media
	req := db.Model(md)
	if limit != 0 {
		req = req.Limit(limit)
	}
	return media, req.Order("id DESC").Find(&media).Error
}

func (md *Media) GetRandom(db *gorm.DB, rnd int32) (*Media, error) {
	media := &Media{}

	if err := db.Model(md).Where("`id`=((? % (select count(*) from `media`))+1)", rnd).First(&media).Error; err != nil {
		return media, err
	} else {
		return media, err
	}
}

func (md *Media) GetTo(db *gorm.DB, last_id int, limit int) ([]Media, error) {
	var media []Media
	return media, db.Model(md).Where("id < ?", last_id).Limit(limit).Order("id DESC").Find(&media).Error
}

func (md *Media) GetByHash(db *gorm.DB, hash string) (*Media, bool, error) {
	err := db.Where(&Media{Hash: hash}).First(md).Error
	if err == gorm.ErrRecordNotFound {
		return md, false, nil
	} else {
		return md, true, err
	}
}

func mediaTypeFromContentType(contenttype string) string {
	return strings.ToLower(strings.SplitN(contenttype, "/", 2)[0])
}

func isContentType(ct string) bool {
	return len(ct) != 0 && strings.ToLower(ct) != "application/octet-stream" && strings.IndexRune(ct, '/') != 0
}

func generateThumb(set *settings.Settings, fname, ctype string) (th *string, err error) {
	switch ctype {
	case "video":
		videothumb := fname + ".png"
		cmd := exec.Command(set.FFmpeg, "-i", path.Join(set.DataPath, fname),
			"-vf", `scale=w='min(1\,min(640/iw\,360/ih))*640':h=-1`,
			"-vframes", "1", path.Join(set.DataPath, videothumb))

		return &videothumb, cmd.Run()
	}
	return
}

var nextTempFileId uint32 = rand.Uint32()

func (md *Media) NewFromFile(db *gorm.DB, rf io.Reader, contenttype string, set *settings.Settings) (*Media, error) {
	tempFileName := path.Join(path.Dir(set.DataPath), fmt.Sprintf("_tmp_%v.tmp", atomic.AddUint32(&nextTempFileId, 1)))

	of, err := os.Create(tempFileName)
	if err != nil {
		return nil, err
	}

	tempneedclose := true
	tempneedremove := true
	defer func() {
		if tempneedclose {
			of.Close()
		}
		if tempneedremove {
			err := os.Remove(tempFileName)
			if err != nil {
				log.Log.Errorf("Error removing tempfile: %v", err)
			}
		}
	}()

	var fsize int64
	var buffer []byte = make([]byte, 1024*128) // 128 kbytes
	var readerr error
	var readcount int
	hmd5 := md5.New()

	for step := 0; readerr != io.EOF; step++ {
		readcount, readerr = rf.Read(buffer)

		if readerr != nil && readerr != io.EOF {
			return nil, err
		}

		if step == 0 {
			detected_ctype := http.DetectContentType(buffer[:readcount])
			if isContentType(detected_ctype) {
				contenttype = detected_ctype
			}
		}

		fsize += int64(readcount)

		_, err = hmd5.Write(buffer[:readcount])
		if err != nil {
			return nil, err
		}

		_, wrerr := of.Write(buffer[:readcount])
		if wrerr != nil {
			return nil, wrerr
		}
	}

	tempneedclose = false
	of.Close()

	mediatype := mediaTypeFromContentType(contenttype)

	switch mediatype {
	case "video", "audio", "image":
	default:
		log.Log.Infof("Aborted content type %v", contenttype)
		return nil, ErrIncorrectContentType
	}

	hash := base64.URLEncoding.EncodeToString(hmd5.Sum(nil))[0:22] // len(base64(md5)) == 22

	model, exists, err := md.GetByHash(db, hash)
	if err != nil {
		return nil, fmt.Errorf("Media selecting gorm error: %v\n", err)
	}

	if exists {
		return model, ErrAlreadyUploaded
	}

	dbfilename := path.Join(hash[:1], fmt.Sprintf("%c_%s", mediatype[0], hash[1:]))
	filename := path.Join(path.Dir(set.DataPath), dbfilename)

	if err = os.MkdirAll(path.Dir(filename), 0666); err != nil {
		return nil, err
	}

	if err = os.Rename(tempFileName, filename); err != nil {
		return nil, err
	}
	tempneedremove = false

	thumb, err := generateThumb(set, dbfilename, mediatype)
	if err != nil {
		log.Log.Errorf("Error generating thumb for file %s: %v", filename, err)
		return nil, err
	}

	model, err = model.New(db, dbfilename, hash, contenttype, fsize, thumb)
	if err != nil {
		remerr := os.Remove(filename)
		if remerr != nil {
			log.Log.Errorf("Error removing not media file %s: %v", filename, remerr)
		}
		if thumb != nil {
			remerr = os.Remove(*thumb)
			if remerr != nil {
				log.Log.Errorf("Error removing thumb file %s: %v", *thumb, remerr)
			}
		}
		return nil, fmt.Errorf("Media creating gorm error: %v\n", err)
	}
	return model, nil
}

func (md *Media) GenerateThumbnail(db *gorm.DB, set *settings.Settings) (*string, error) {
	// remove old thumbnail
	if md.Thumbnail != nil {
		oldpath := path.Join(set.DataPath, *md.Thumbnail)
		if _, err := os.Stat(oldpath); err == nil || os.IsExist(err) {
			err = os.Remove(oldpath)
			if err != nil {
				return nil, err
			}
		}
	}

	// generate new
	return generateThumb(set, md.File, mediaTypeFromContentType(md.Type))
}

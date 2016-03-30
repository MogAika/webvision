package handlers

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"

	"github.com/jinzhu/gorm"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/settings"
	"github.com/mogaika/webvision/views"
)

var nextTempFileId uint32 = rand.Uint32()

var ErrAlreadyUploaded = errors.New("File already uploaded on server")
var ErrIncorrectContentType = errors.New("Server support only image/video/audio files")

func isContentType(ct string) bool {
	return len(ct) != 0 && strings.ToLower(ct) != "application/octet-stream" && strings.IndexRune(ct, '/') != 0
}

func ProcessFile(db *gorm.DB, rf multipart.File, contenttype string, set *settings.Settings) (*models.Media, error) {
	tempFileName := path.Join(path.Dir(set.DataPath), fmt.Sprintf("_tmp_%v.tmp", atomic.AddUint32(&nextTempFileId, 1)))

	of, err := os.Create(tempFileName)
	if err != nil {
		return nil, err
	}
	defer of.Close()

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
			if !isContentType(detected_ctype) {
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
	of.Close()

	mediatype := strings.ToLower(strings.SplitN(contenttype, "/", 2)[0])

	switch mediatype {
	case "video", "audio", "image":
	default:
		return nil, ErrIncorrectContentType
	}

	hash := base64.URLEncoding.EncodeToString(hmd5.Sum(nil))[0:22] // len(md5) == 22

	model, exists, err := (&models.Media{}).GetByHash(db, hash)
	if err != nil {
		return nil, fmt.Errorf("Media selecting gorm error: %v\n", err)
	}

	if exists {
		return model, ErrAlreadyUploaded
	}

	filename := path.Join(path.Dir(set.DataPath), hash[:1], fmt.Sprintf("%c_%s", mediatype[0], hash[1:]))

	if err = os.MkdirAll(path.Dir(filename), 0666); err != nil {
		return nil, err
	}

	if err = os.Rename(tempFileName, filename); err != nil {
		return nil, err
	}

	var thumb *string = nil

	model, err = model.New(db, filename, hash, contenttype, fsize, thumb)
	if err != nil {
		return nil, fmt.Errorf("Media creating gorm error: %v\n", err)
	}
	return model, nil
}

func HandlerUploadGet(w http.ResponseWriter, r *http.Request) {
	views.ViewUpload(w)
}

func HandlerUploadPost(w http.ResponseWriter, r *http.Request) {
	db, set := VarsFromRequest(r)
	r.Body = http.MaxBytesReader(w, r.Body, set.MaxDataSize)

	f, fh, err := r.FormFile("heh")

	if err == nil {
		defer f.Close()
		_, err := ProcessFile(db, f, fh.Header.Get("Content-Type"), set)
		if err == nil {
			views.ViewUploadResult(w, "")
		} else {
			views.ViewUploadResult(w, err.Error())
		}
	} else {
		views.ViewError(w, 500, "Error", err.Error())
	}
}

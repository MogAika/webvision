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
	"github.com/mogaika/webvision/log"
	"github.com/mogaika/webvision/models"
	"github.com/mogaika/webvision/settings"
	"github.com/mogaika/webvision/views"
)

var nextTempFileId uint32 = rand.Uint32()

var ErrAlreadyUploaded = errors.New("File alrady uploaded on server")

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

		if step == 0 && !isContentType(contenttype) {
			contenttype = http.DetectContentType(buffer[:readcount])
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

	hash := base64.URLEncoding.EncodeToString(hmd5.Sum(nil))[0:22] // len(md5) == 22
	log.Log.Info(len(hash), hash)

	model, err := (&models.Media{}).GetByHash(db, hash)
	if err != nil {
		return nil, err
	}

	if model.ID != 0 {
		return model, ErrAlreadyUploaded
	}

	filename := path.Join(path.Dir(set.DataPath), hash[:1], hash[1:])

	if err = os.MkdirAll(path.Dir(filename), 0666); err != nil {
		return nil, err
	}

	if err = os.Rename(tempFileName, filename); err != nil {
		return nil, err
	}

	return model.New(db, filename, hash, contenttype, fsize, nil)
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
		model, err := ProcessFile(db, f, fh.Header.Get("Content-Type"), set)
		if err == nil {
			views.ViewError(w, 200, "Uploaded", fh.Filename+" uploaded. Type: "+model.Type)
		} else if err == ErrAlreadyUploaded {
			views.ViewError(w, 200, "File already uploaded", "")
		} else {
			log.Log.Errorf("Error processing file \"%s\": %v", fh.Filename, err)
			views.ViewError(w, 500, "Error", "Error processing file on server side")
		}
	} else {
		views.ViewError(w, 500, "Error", err.Error())
	}
}

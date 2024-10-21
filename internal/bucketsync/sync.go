package bucketsync

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"ossmark/internal/config"
	"ossmark/pkg/bucketutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type fileStat struct {
	Path       string
	Size       int
	ModifiedTm time.Time
}

func Sync(conf *config.BucketConfig, dir, mode string) error {
	b, err := bucketutil.NewBucket(&bucketutil.AccessKey{Id: conf.AccessKeyId, Secret: conf.AcessKeySecret}, conf.BucketName, conf.BucketLocation)
	if err != nil {
		return errors.Wrap(err, "NewBucket")
	}

	var workdir string
	if dir == "" {
		workdir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
		workdir = path.Join(workdir, ".ossmark-data", b.BucketName)
	} else {
		workdir = dir
	}
	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return err
	}

	err = syncLocal(b, workdir, mode)
	if err != nil {
		return err
	}
	return syncRemote(b, workdir, mode)
}

func syncLocal(b *oss.Bucket, workdir, mode string) error {
	logrus.WithFields(logrus.Fields{"work_dir": workdir, "bucket": b.BucketName, "mode": mode}).Info("Sync local")

	handler := func(s *fileStat) error {
		if s.Size == 0 {
			return nil
		}

		f, err := os.Open(s.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		key := strings.TrimPrefix(s.Path, workdir+"/")
		l := logrus.WithField("key", key)

		resp, err := b.GetObjectMeta(key)
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchKey") {
				l.Info("Put object to remote")
				return b.PutObject(key, f)
			}
			return err
		}
		etag := strings.Trim(resp.Get("Etag"), "\"")
		modified := resp.Get("Last-Modified")

		sign, err := md5File(f)
		if err != nil {
			return err
		}
		f.Seek(0, 0)

		l.WithFields(logrus.Fields{"local": sign, "remote": etag}).Debug("Check etag")
		if etag == sign {
			return nil
		}

		timeCheck := false
		if mode == "time" {
			modifiedTm, err := time.Parse(time.RFC1123, modified)
			if err != nil {
				return nil
			}
			l.WithFields(logrus.Fields{"local": formatTm(s.ModifiedTm), "remote": formatTm(modifiedTm)}).Debug("Check last modified time")
			if s.ModifiedTm.After(modifiedTm) {
				return nil
			}
			timeCheck = true
		}

		if mode == "local" || timeCheck {
			l.Info("Put object to remote")
			return b.PutObject(key, f)
		}
		return nil
	}

	return filepath.Walk(workdir, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		return handler(&fileStat{Path: p, Size: int(info.Size()), ModifiedTm: info.ModTime()})
	})
}

func syncRemote(b *oss.Bucket, workdir, mode string) error {
	logrus.WithFields(logrus.Fields{"work_dir": workdir, "bucket": b.BucketName, "mode": mode}).Info("Sync remote")

	return bucketutil.ListObjectsWithHandler(b, func(obj *oss.ObjectProperties) error {
		l := logrus.WithField("key", obj.Key)

		dir := path.Dir(obj.Key)
		err := os.MkdirAll(path.Join(workdir, dir), 0755)
		if err != nil {
			return err
		}

		name := path.Join(workdir, dir, path.Base(obj.Key))
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		sign, err := md5File(f)
		if err != nil {
			return err
		}
		etag := strings.Trim(obj.ETag, "\"")
		l.WithFields(logrus.Fields{"local": sign, "remote": etag}).Debug("Check etag")
		if etag == sign {
			return nil
		}

		timeCheck := false
		if mode == "time" {
			st, err := os.Stat(name)
			if err != nil {
				return err
			}
			l.WithFields(logrus.Fields{"local": formatTm(st.ModTime()), "remote": formatTm(obj.LastModified)}).Debug("Check last modified time")
			if st.ModTime().Before(obj.LastModified) {
				return nil
			}
			timeCheck = true
		}

		if mode == "remote" || timeCheck {
			l.WithFields(logrus.Fields{"key": obj.Key, "etag": etag, "last_modified": obj.LastModified}).Info("Get object from remote")
			return b.GetObjectToFile(obj.Key, name)
		}
		return nil
	})
}

func md5sum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func md5File(file *os.File) (string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(md5sum(content)), nil
}

func formatTm(t time.Time) string {
	return t.Local().Format("2006/01/02 15:04:05")
}

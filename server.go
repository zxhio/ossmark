package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
)

var (
	cryptoAESKey [32]byte
)

func init() {
	_, err := rand.Read(cryptoAESKey[:])
	if err != nil {
		panic(err)
	}
}

func serve(b *oss.Bucket) {
	logrus.WithField("bucket", b.BucketName).Info("Start article server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		months := make(map[string]monthArticleList)
		err := listObjects(b, func(obj *oss.ObjectProperties) error {
			if strings.HasSuffix(obj.Key, "/") {
				return nil
			}

			modifyTm := obj.LastModified.Local().Format("2006/01/02 15:04:05")
			month := obj.LastModified.Local().Format("2006/01")

			encryptedObjKey, err := EncryptAES_CBC([]byte(obj.Key), aesRandKey)
			if err != nil {
				return err
			}

			m := months[month]
			m.Month = month
			m.Articles = append(m.Articles, article{
				Link:       fmt.Sprintf("%s?key=%s&modify_tm=%s", path.Join("articles", path.Base(obj.Key)), url.QueryEscape(encryptedObjKey), url.QueryEscape(modifyTm)),
				Name:       strings.TrimSuffix(path.Base(obj.Key), ".md"),
				LastModify: modifyTm,
			})
			months[m.Month] = m
			return nil
		})
		if err != nil {
			fmt.Fprintf(w, "encryt error %v", err)
			return
		}

		list := make([]monthArticleList, len(months))
		for _, v := range months {
			sort.Sort(ArticleSlice(v.Articles))
			list = append(list, v)
		}
		sort.Sort(MonthArticleSlice(list))
		ArticleListTmpl.Execute(w, articleListBody{MonthArticles: list})
	})

	http.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		queryKey := r.URL.Query().Get("key")
		key, err := DecryptAES_CBC(queryKey, aesRandKey)
		if err != nil {
			fmt.Fprintf(w, "internal error")
			return
		}

		reader, err := b.GetObject(key)
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchKey") {
				w.Write([]byte("No such article\n"))
			} else {
				w.Write([]byte(err.Error()))
			}
			return
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		output := blackfriday.Run(content, blackfriday.WithExtensions(blackfriday.CommonExtensions))
		err = ArticleContentTmpl.Execute(w, articleContentBody{
			Title:    strings.TrimSuffix(path.Base(key), ".md"),
			Content:  string(output),
			ModifyTm: r.URL.Query().Get("modify_tm"),
		})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	})

	if ossmarkConf.ListenPort == 0 {
		ossmarkConf.ListenPort = 9991
	}
	l, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", ossmarkConf.ListenPort))
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to net.Listen")
	}
	logrus.WithField("addr", l.Addr()).Info("Listen on")

	http.Serve(l, nil)
}

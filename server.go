package main

import (
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

func serve(b *oss.Bucket) {
	logrus.WithField("bucket", b.BucketName).Info("Start article server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		months := make(map[string]monthArticleList)
		listObjects(b, func(obj *oss.ObjectProperties) error {
			if strings.HasSuffix(obj.Key, "/") {
				return nil
			}

			modifyTm := obj.LastModified.Format("2006/01/02 15:04:05")
			month := obj.LastModified.Format("2006/01")

			m := months[month]
			m.Month = month
			m.Articles = append(m.Articles, article{
				Link:       fmt.Sprintf("%s?modify_tm=%s", path.Join("articles", obj.Key), url.QueryEscape(modifyTm)),
				Name:       strings.TrimSuffix(path.Base(obj.Key), ".md"),
				LastModify: modifyTm,
			})
			months[m.Month] = m
			return nil
		})
		list := make([]monthArticleList, len(months))
		for _, v := range months {
			list = append(list, v)
		}
		sort.Sort(monthArticles(list))
		ArticleListTmpl.Execute(w, articleListBody{MonthArticles: list})
	})

	http.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		key := strings.TrimPrefix(r.URL.String(), "/articles/")
		key, err := url.PathUnescape(key)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		s1, _, _ := strings.Cut(key, "?")
		key = s1

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

	if Conf.ListenPort == 0 {
		Conf.ListenPort = 9991
	}
	l, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", Conf.ListenPort))
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to net.Listen")
	}
	logrus.WithField("addr", l.Addr()).Info("Listen on")

	http.Serve(l, nil)
}

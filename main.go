package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type AccessKey struct {
	Id     string `json:"access_key_id"`
	Secret string `json:"access_key_secret"`
}

type Config struct {
	AccessKey
	BucketName      string `json:"bucket_name"`
	SkipObjectRegex string `json:"skip_object_regex"`
	ListenPort      int    `json:"listen_port"`
}

func readAndParseConfig(confPath string) (*Config, error) {
	content, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(content, &c)
	return &c, err
}

func main() {
	confPath := flag.StringP("conf", "c", "/etc/ossmark.json", "config path")
	flag.Parse()

	conf, err := readAndParseConfig(*confPath)
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to read config file")
	}

	// logrus.SetOutput(&lumberjack.Logger{
	// 	Filename:   "/var/log/ossmark.log",
	// 	MaxSize:    32, // in MB
	// 	MaxBackups: 10,
	// 	Compress:   true,
	// })

	articleListTmpl, err := template.New("article_list").Parse(articleListTemplate)
	if err != nil {
		logrus.WithError(err).Fatal("Fail to new article_list template")
	}
	articleContentTmpl, err := template.New("article_content").Parse(articleContentTemplate)
	if err != nil {
		logrus.WithError(err).Fatal("Fail to new article_content template")
	}

	b, err := newBucket(&conf.AccessKey, conf.BucketName, "oss-cn-hangzhou")
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to new bucket client")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		months := make(map[string]monthArticleList)
		listObjects(b, func(obj *oss.ObjectProperties) error {
			month := obj.LastModified.Format("2006/01")
			m := months[month]
			m.Month = month
			m.Articles = append(m.Articles, article{
				Path:       path.Join("articles", obj.Key),
				Name:       path.Base(obj.Key),
				LastModify: obj.LastModified.Format(time.RFC3339),
			})
			months[m.Month] = m
			return nil
		})
		list := make([]monthArticleList, len(months))
		for _, v := range months {
			list = append(list, v)
		}
		sort.Sort(monthArticles(list))
		articleListTmpl.Execute(w, articleListBody{MonthArticles: list})
	})

	http.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

		key := strings.TrimPrefix(r.URL.String(), "/articles/")
		key, err := url.QueryUnescape(key)
		if err != nil {
			w.Write([]byte(err.Error()))
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
		err = articleContentTmpl.Execute(w, articleContentBody{Title: path.Base(key), Content: string(output)})
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	})

	if conf.ListenPort == 0 {
		conf.ListenPort = 9991
	}
	l, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", conf.ListenPort))
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to net.Listen")
	}
	logrus.WithField("addr", l.Addr()).Info("Listen on")

	http.Serve(l, nil)
}

type ObjectHandler func(obj *oss.ObjectProperties) error

func listObjects(b *oss.Bucket, handler ObjectHandler) error {
	var (
		makKeys   = 100
		nextToken string
	)
	for {
		resp, err := b.ListObjectsV2(oss.MaxKeys(makKeys), oss.ContinuationToken(nextToken))
		if err != nil {
			return err
		}
		for _, o := range resp.Objects {
			if handler != nil {
				err = handler(&o)
				if err != nil {
					return err
				}
			}
		}
		if len(resp.Objects) < makKeys || resp.NextContinuationToken == "" {
			return err
		}
		nextToken = resp.NextContinuationToken
	}
}

func newBucket(ak *AccessKey, bucketName, location string) (*oss.Bucket, error) {
	client, err := oss.New(fmt.Sprintf("%s.aliyuncs.com", location), ak.Id, ak.Secret)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketName)
}

// func newBucketWithIntranet(ak *AccessKey, bucketName, location string) (*oss.Bucket, error) {
// 	client, err := oss.New(fmt.Sprintf("%s-internal.aliyuncs.com", location), ak.Id, ak.Secret)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.Bucket(bucketName)
// }

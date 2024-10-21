package article

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"ossmark/internal/config"
	"ossmark/pkg/bucketutil"
	"ossmark/pkg/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
)

type ArticleServer struct {
	bucket *oss.Bucket
}

func NewArticleServer(conf *config.BucketConfig) (*ArticleServer, error) {
	b, err := bucketutil.NewBucket(&bucketutil.AccessKey{Id: conf.AccessKeyId, Secret: conf.AcessKeySecret}, conf.BucketName, conf.BucketLocation)
	if err != nil {
		return nil, errors.Wrap(err, "NewBucket")
	}
	return &ArticleServer{bucket: b}, nil
}

func (s *ArticleServer) BucketName() string {
	return s.bucket.BucketName
}

func (s *ArticleServer) Serve(lis net.Listener) error {
	http.HandleFunc("/", s.listArticle)
	http.HandleFunc("/articles/", s.showArticle)

	return http.Serve(lis, nil)
}

func (s *ArticleServer) listArticle(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

	months := make(map[string]monthArticleList)
	err := bucketutil.ListObjectsWithHandler(s.bucket, func(obj *oss.ObjectProperties) error {
		if strings.HasSuffix(obj.Key, "/") {
			return nil
		}

		modifyTm := obj.LastModified.Local().Format("2006/01/02 15:04:05")
		month := obj.LastModified.Local().Format("2006/01")

		encryptedObjKey, err := utils.EncryptAES_CBC([]byte(obj.Key), utils.RandAESKey)
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
}

func (s *ArticleServer) showArticle(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{"addr": r.RemoteAddr, "url": r.URL}).Info("New connection")

	queryKey := r.URL.Query().Get("key")
	key, err := utils.DecryptAES_CBC(queryKey, utils.RandAESKey)
	if err != nil {
		fmt.Fprintf(w, "internal error")
		return
	}

	reader, err := s.bucket.GetObject(key)
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
}

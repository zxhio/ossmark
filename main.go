package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/natefinch/lumberjack"
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
	WorkDir         string `json:"work_dir"`
}

var (
	Conf Config
)

func readAndParseConfig(confPath string) error {
	content, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &Conf)
}

type syncFlag struct {
	set   bool
	value string
}

func (f *syncFlag) String() string { return f.value }
func (f *syncFlag) Set(s string) error {
	f.value = s
	f.set = true
	return nil
}
func (f *syncFlag) Type() string { return "bool|string" }

func main() {
	var sf syncFlag
	f := flag.CommandLine.VarPF(&sf, "sync", "", "sync bucket base on [time|local|remote], default 'time'")
	f.NoOptDefVal = ""

	confPath := flag.String("conf", "conf/ossmark.json", "config path")
	enableArticle := flag.Bool("article", false, "start a server to show article")
	flag.Parse()

	err := readAndParseConfig(*confPath)
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to read config file")
	}
	logrus.SetLevel(logrus.DebugLevel)

	b, err := newBucket(&Conf.AccessKey, Conf.BucketName, "oss-cn-hangzhou")
	if err != nil {
		logrus.WithError(err).Fatal("Fatal to new bucket client")
	}

	if sf.set {
		err = sync(b, Conf.WorkDir, sf.value)
		if err != nil {
			logrus.WithError(err).Fatal("Fatal to sync bucket")
		}
		return
	}

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/ossmark.log",
		MaxSize:    32, // in MB
		MaxBackups: 10,
		Compress:   true,
	})
	if *enableArticle {
		serve(b)
	}
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

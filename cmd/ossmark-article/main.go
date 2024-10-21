package main

import (
	"fmt"
	"net"

	"ossmark/internal/article"
	"ossmark/internal/config"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type Config struct {
	*config.BucketConfig
	ListenPort int    `json:"listen_port"`
	LogPath    string `json:"log_path"`
}

func main() {
	confPath := pflag.String("config", "", "bucket config path")
	listenPort := pflag.Int("listen_port", 0, "article server port")
	logPath := pflag.String("log_path", "", "article server log path")
	config.SetBucketFlags()
	pflag.Parse()

	var conf *Config
	if *confPath != "" {
		var err error
		conf, err = config.New[Config](*confPath)
		if err != nil {
			panic(err)
		}

		// if command not specify, use config fields
		if *listenPort == 0 {
			*listenPort = conf.ListenPort
		}
		if *logPath == "" {
			*logPath = "logs/ossmark-article.log"
		}
	} else {
		conf = &Config{BucketConfig: config.GetConfigByFlags()}
	}

	if *listenPort == 0 {
		*listenPort = 9346
	}

	logrus.SetLevel(logrus.DebugLevel)
	if *logPath != "" {
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   *logPath,
			MaxSize:    32, // in MB
			MaxBackups: 10,
			Compress:   true,
		})
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("[::]:%d", *listenPort))
	if err != nil {
		panic(err)
	}
	logrus.WithField("addr", lis.Addr()).Info("Listen on")

	s, err := article.NewArticleServer(conf.BucketConfig)
	if err != nil {
		panic(err)
	}

	logrus.WithField("name", s.BucketName()).Info("Serve bucket article")
	err = s.Serve(lis)
	panic(err)
}

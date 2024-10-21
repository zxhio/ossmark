package main

import (
	"ossmark/internal/bucketsync"
	"ossmark/internal/config"

	"github.com/spf13/pflag"
)

type Config struct {
	*config.BucketConfig
	WorkDir string `json:"work_dir"`
	Mode    string `json:"mode"`
}

func main() {
	confPath := pflag.String("config", "", "bucket config path")
	mode := pflag.String("mode", "", "oss bucket sync mode base on [ time | local | remote ]")
	workDir := pflag.String("work_dir", "", "oss bucket sync work dir")
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
		if *mode == "" {
			*mode = conf.Mode
		}
		if *workDir == "" {
			*workDir = "data/sync-data"
		}
	} else {
		conf = &Config{BucketConfig: config.GetConfigByFlags()}
	}

	if *mode == "" {
		*mode = "time"
	}

	err := bucketsync.Sync(conf.BucketConfig, *workDir, *mode)
	if err != nil {
		panic(err)
	}
}

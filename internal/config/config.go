package config

import (
	"encoding/json"
	"os"

	"github.com/spf13/pflag"
)

type BucketConfig struct {
	BucketName     string `json:"bucket_name"`
	BucketLocation string `json:"bucket_location"`
	AccessKeyId    string `json:"access_key_id"`
	AcessKeySecret string `json:"access_key_secret"`
}

var (
	bucketName      *string
	bucketLocation  *string
	accessKeyId     *string
	accessKeySecret *string
)

func SetBucketFlags() {
	bucketName = pflag.String("bucket_name", "", "bucket name")
	bucketLocation = pflag.String("bucket_location", "", "bucket location, eg. oss-cn-hangzhou")
	accessKeyId = pflag.String("access_key_id", "", "access key id")
	accessKeySecret = pflag.String("access_key_secret", "", "access key secret")
}

func GetConfigByFlags() *BucketConfig {
	return &BucketConfig{
		BucketName:     *bucketName,
		BucketLocation: *bucketLocation,
		AccessKeyId:    *accessKeyId,
		AcessKeySecret: *accessKeySecret,
	}
}

func New[T any](confPath string) (*T, error) {
	content, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	var c T
	err = json.Unmarshal(content, &c)
	return &c, err
}

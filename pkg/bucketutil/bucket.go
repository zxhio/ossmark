package bucketutil

import (
	"fmt"
	"net"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AccessKey struct {
	Id     string
	Secret string
}

func NewBucket(ak *AccessKey, bucketName, location string) (*oss.Bucket, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s-internal.aliyuncs.com:443", location), time.Second)
	if err != nil {
		return NewBucketWithPublic(ak, bucketName, location)
	}
	defer conn.Close()

	return NewBucketWithIntranet(ak, bucketName, location)
}

// 外网访问
func NewBucketWithPublic(ak *AccessKey, bucketName, location string) (*oss.Bucket, error) {
	client, err := oss.New(fmt.Sprintf("%s.aliyuncs.com", location), ak.Id, ak.Secret)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketName)
}

// 内网访问
func NewBucketWithIntranet(ak *AccessKey, bucketName, location string) (*oss.Bucket, error) {
	client, err := oss.New(fmt.Sprintf("%s-internal.aliyuncs.com", location), ak.Id, ak.Secret)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketName)
}

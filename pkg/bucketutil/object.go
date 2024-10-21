package bucketutil

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

type ObjectHandler func(obj *oss.ObjectProperties) error

func ListObjectsWithHandler(b *oss.Bucket, handler ObjectHandler) error {
	var (
		makKeys   = 100
		nextToken string
	)

	for {
		resp, err := b.ListObjectsV2(oss.MaxKeys(makKeys), oss.ContinuationToken(nextToken))
		if err != nil {
			return errors.Wrap(err, "ListObjectsV2")
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
			break
		}
		nextToken = resp.NextContinuationToken
	}

	return nil
}

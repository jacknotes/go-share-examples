package aliyun

import (
	"fmt"
	"jacknotes/go-share-examples/cloud-station/pro/store"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func NewAliyunOssUpload(endpoint, ak, sk string) store.OSSUploader {
	return &impl{
		endpoint: endpoint,
		ak:       ak,
		sk:       sk,
	}
}

// 这个对象，实现我们定义的接口
type impl struct {
	// oss 服务
	endpoint string
	ak       string
	sk       string
}

func (i *impl) Upload(bucketName, objectKey, fileName string) (downloadURL string, err error) {
	// 连接OSS
	client, err := oss.New(i.endpoint, i.ak, i.sk)
	if err != nil {
		err = fmt.Errorf("new client error, %s\n", err)
		return
	}
	// 切换bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		err = fmt.Errorf("get bucket %s error, %s\n", bucketName, err)
		return
	}
	// 上传文件,arg1=oss存储的文件名称， arg2=原文件地址
	err = bucket.PutObjectFromFile(objectKey, fileName)
	if err != nil {
		err = fmt.Errorf("upload file  %s error, %s\n", fileName, err)
		return
	}

	// 生成下载链接
	return bucket.SignURL(objectKey, oss.HTTPGet, 60*60*24*3)
}

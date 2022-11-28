package aliyun

import (
	"fmt"
	"jacknotes/go-share-examples/cloud-station/pro/store"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10" //validator v10版本
)

var (
	// 新建校验器
	validate = validator.New()
)

func NewAliyunOssUpload(endpoint, ak, sk string) (store.OSSUploader, error) {
	uploader := &impl{
		Endpoint: endpoint,
		Ak:       ak,
		Sk:       sk,
	}

	if err := uploader.Validate(); err != nil {
		return nil, fmt.Errorf("validate params error %s", err)
	}

	return uploader, nil
}

// 这个对象，实现我们定义的接口
type impl struct {
	// oss 服务
	Endpoint string `validate:"required"` //校验参数是需要的，使用反射时，属性需要大写开头
	Ak       string `validate:"required"`
	Sk       string `validate:"required"`
}

func (i *impl) Validate() error {
	// 校验结构体impl中的参数 是否校验通过
	return validate.Struct(i)
}

func (i *impl) Upload(bucketName, objectKey, fileName string) (downloadURL string, err error) {
	// 连接OSS
	client, err := oss.New(i.Endpoint, i.Ak, i.Sk)
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

package main

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	conf     = NewDefaultConfig()
	fileName = "go.mod"
)

type Config struct {
	Endpoint   string
	AK         string
	SK         string
	BucketName string
}

func NewDefaultConfig() *Config {
	return &Config{
		BucketName: "devopscloud-station",
	}
}

func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("oss endpoint required, use env ALI_OSS_ENDPOINT to set")
	}
	if c.AK == "" {
		return fmt.Errorf("ali AccessKey required, use env ALI_AK to set")
	}
	if c.SK == "" {
		return fmt.Errorf("ali SecretKey required, use env ALI_SK to set")
	}
	return nil
}

func LoadConfigFromEnv() {
	conf.Endpoint = os.Getenv("ALI_OSS_ENDPOINT")
	conf.AK = os.Getenv("ALI_AK")
	conf.SK = os.Getenv("ALI_SK")
}

func UploadFile(filename string) error {
	// 连接OSS
	client, err := oss.New(conf.Endpoint, conf.AK, conf.SK)
	if err != nil {
		return fmt.Errorf("new client error, %s\n", err)
	}
	// 切换bucket
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		return fmt.Errorf("get bucket %s error, %s\n", conf.BucketName, err)
	}
	// 上传文件,arg1=oss存储的文件名称， arg2=原文件地址
	err = bucket.PutObjectFromFile(filename, filename)
	if err != nil {
		return fmt.Errorf("upload file to bucket %s error, %s\n", conf.BucketName, err)
	}

	return nil
}

// main是拼装流程的
func main() {
	// 加载配置
	LoadConfigFromEnv()

	// 校验配置
	if err := conf.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 上传文件
	if err := UploadFile(fileName); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// 正常退出
	fmt.Printf("文件: %s 上传成功", fileName)
	os.Exit(0)
}

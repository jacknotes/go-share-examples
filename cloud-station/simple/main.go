package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	conf = NewDefaultConfig()
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

func UploadFile(filename string) (downloadURL string, err error) {
	// 连接OSS
	client, err := oss.New(conf.Endpoint, conf.AK, conf.SK)
	if err != nil {
		err = fmt.Errorf("new client error, %s\n", err)
		return
	}
	// 切换bucket
	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		err = fmt.Errorf("get bucket %s error, %s\n", conf.BucketName, err)
		return
	}
	// 上传文件,arg1=oss存储的文件名称， arg2=原文件地址
	err = bucket.PutObjectFromFile(filename, filename)
	if err != nil {
		err = fmt.Errorf("upload file to bucket %s error, %s\n", conf.BucketName, err)
		return
	}

	// 生成下载链接
	return bucket.SignURL(filename, oss.HTTPGet, 60*60*24*3)
}

/*
	CLI 说明
*/

var (
	fileName string
	help     bool
)

// 声明CLI的参数
func init() {
	// 声明cli参数
	flag.StringVar(&fileName, "f", "", "请输入需要上传的文件路径")
	flag.BoolVar(&help, "h", false, "打印本工具的使用说明")
}

// 工具使用说明函数
func usage() {
	fmt.Fprintf(os.Stderr, `cloud-station version: 0.0.1
Usage: cloud-station [-h] -f <uplaod_file_path>
Options:
`)
	// 打印默认参数
	flag.PrintDefaults()
}

// 负责接收用户传入的参数
func LoadArgsFromCLI() {
	// 解析CLI参数，并赋值给变量
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}
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

	// 接收用户参数
	LoadArgsFromCLI()

	// 上传文件
	downloadURL, err := UploadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// 正常退出
	fmt.Printf("文件: %s 上传成功\n", fileName)

	// 打印下载链接
	fmt.Printf("下载链接: %s\n", downloadURL)
	fmt.Println("注意: 文件下载有效期为3天, 中转站保存时间为3天, 请及时下载")

	os.Exit(0)
}

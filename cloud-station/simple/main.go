package main

import (
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	conf = NewDefaultConfig()
)

func NewDefaultConfig() *Config {
	return &Config{
		BucketName: "devopscloud-station",
	}
}

type Config struct {
	Endpoint   string
	AK         string
	SK         string
	BucketName string
}

func LoadConfigFromEnv() {
	conf.Endpoint = os.Getenv("ALI_OSS_ENDPOINT")
	conf.AK = os.Getenv("ALI_AK")
	conf.SK = os.Getenv("ALI_SK")
}

func main() {
	LoadConfigFromEnv()

	client, err := oss.New(conf.Endpoint, conf.AK, conf.SK)
	if err != nil {
		panic(err)
	}

	bucket, err := client.Bucket(conf.BucketName)
	if err != nil {
		panic(err)
	}

	err = bucket.PutObjectFromFile("my_go_mod", "go.mod")
	if err != nil {
		panic(err)
	}
}

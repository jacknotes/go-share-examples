package aliyun_test

import (
	"jacknotes/go-share-examples/cloud-station/pro/store/aliyun"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ep string
	ak string
	sk string
	bn string
)

// TDD 测试驱动开发
func TestUpload(t *testing.T) {
	// 新建一个断言对象
	should := assert.New(t)

	uploader, err := aliyun.NewAliyunOssUpload(ep, ak, sk)
	if should.NoError(err) {
		// 测试aliyun impl upload方法
		downloadUrl, err := uploader.Upload(bn, "impl.go", "impl.go")

		// //当前路径
		// workdir, _ := os.Getwd()
		// fmt.Println("workdir:", workdir)

		// // 为了简化逻辑，可以更形象的表现出来，引用后面的assert
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// if downloadUrl == "" {
		// 	t.Fatal("no download")
		// }

		// if err == nil，相当于should.NoError(err)，如果没有error
		if should.NoError(err) {
			// downloadURL不为空
			should.NotEmpty(downloadUrl)
		}
	}

}

// 测试时候 通过环境变量加载参数
func init() {
	ep = os.Getenv("ALI_OSS_ENDPOINT")
	ak = os.Getenv("ALI_AK")
	sk = os.Getenv("ALI_SK")
	bn = os.Getenv("ALI_BUCKET_NAME")
}

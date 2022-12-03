package cmd

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2" // 开源好用的第三方库，作用是隐藏输入的密码
	"github.com/spf13/cobra"

	"jacknotes/go-share-examples/cloud-station/pro/store"
	"jacknotes/go-share-examples/cloud-station/pro/store/aliyun"
)

var (
	filename   string
	bucketName string
)

// uploadCmd represents the start command
var uploadCmd = &cobra.Command{
	// upload为子命令名称
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 先初始化secret key
		getSecretKeyFromInputV2()
		// 获取 provider 厂商
		p, err := getProvider()
		if err != nil {
			return err
		}

		if filename == "" {
			return fmt.Errorf("file name required")
		}

		// 为了防止不同用户同一时间上传相同的文件
		// 我们采用用户的主机名作为前置
		hn, err := os.Hostname()
		if err != nil {
			// 如果主机名获取不到，我们使用IP地址
			ipAddr := getOutBindIp()
			if ipAddr == "" {
				hn = "unknown"
			} else {
				hn = ipAddr
			}
		}

		// 为了防止文件都堆在一个文件夹里面 无法查看
		// 我们采用日期进行编码
		day := time.Now().Format("20060102")

		//返回路径的最后一个元素，就是只取文件名
		fn := path.Base(filename)
		// 文件夹格式为20221129/hostname/filename
		objectKey := fmt.Sprintf("%s/%s/%s", day, hn, fn)
		downloadURL, err := p.Upload(bucketName, objectKey, filename)
		if err != nil {
			return err
		}

		// 正常退出
		fmt.Printf("文 件: %s 上传成功\n", filename)

		// 打印下载链接
		fmt.Printf("下载链接: %s\n", downloadURL)
		fmt.Println("注意: 文件下载有效期为3天, 中转站保存时间为3天, 请及时下载")

		return nil
	},
}

func getOutBindIp() string {
	// 向百度发送一个临时链接，取出本地出去的IP地址
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	//取出IP地址
	addr := strings.Split(conn.LocalAddr().String(), ":")
	if len(addr) == 0 {
		return ""
	}

	return addr[0]
}

func getProvider() (p store.OSSUploader, err error) {
	switch ossProvider {
	case "aliyun":
		p, err = aliyun.NewAliyunOssUpload(ossEndpoint, aliAccessKey, aliSecretKey)
		return
	case "qcloud":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [aliyun/qcloud]")
	}
}

func getSecretKeyFromInputV2() {
	prompt := &survey.Password{
		Message: "请输入secret key: ",
	}
	// 调用prompt对象，此时会提示用户输入 secret key，用户输入完成后回车，将会把用户输入的值 赋值给aliSecretKey
	survey.AskOne(prompt, &aliSecretKey)
	fmt.Println()
}

// 初始化函数，执行顺序为var, init(), 普通函数
func init() {
	uploadCmd.PersistentFlags().StringVarP(&bucketName, "bucket_name", "b", "devopscloud-station", "上传文件的目录")
	uploadCmd.PersistentFlags().StringVarP(&filename, "file_name", "f", "", "上传文件的名称")
	// 将uploadCmd配置为RootCmd的子命令
	RootCmd.AddCommand(uploadCmd)
}

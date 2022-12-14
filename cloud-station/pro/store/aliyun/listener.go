package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/k0kubun/go-ansi" // 引入第三方库 修复windows终端输出换行问题
	"github.com/schollz/progressbar/v3"
)

type OssProgressListener struct {
	bar *progressbar.ProgressBar
}

func NewOssProgressListener() *OssProgressListener {
	return &OssProgressListener{}
}

func (l *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	// 调试进度条事件返回的信息，ConsumedBytes、TotalBytes、RwBytes、EventType
	// fmt.Println(event)
	switch event.EventType {
	case oss.TransferStartedEvent:
		// l.bar = progressbar.DefaultBytes(
		// 	event.TotalBytes,
		// 	"文件上传中",
		// )

		l.bar = progressbar.NewOptions64(event.TotalBytes,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(30),
			progressbar.OptionSetDescription("开始上传:"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)
	case oss.TransferDataEvent:
		l.bar.Add64(event.RwBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\n上传完成\n")
	case oss.TransferFailedEvent:
		fmt.Printf("\n上传失败\n")
	default:
	}
}

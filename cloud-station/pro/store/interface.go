package store

// oss 存储适配器
// 阿里云 OSS/Tencent OSS/IDC Minio 开源OSS/AWS S3
type OSSUploader interface {
	// bucketName bucket名称
	// objectKey  上传文件的对象名称
	// fileName  上传的源文件路径
	Upload(bucketName, objectKey, fileName string) (downloadURL string, err error)
}

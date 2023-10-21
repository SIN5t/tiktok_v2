package test

import (
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"testing"
)

func TestMinio(t *testing.T) {
	config := viper.Init("minio")
	picBucket := config.GetString("pic_bucket")

	minio.InitMinIo()
	err := minio.UploadFile(picBucket, "img.png", "./")
	if err != nil {
		klog.Info(err.Error())
	}

}
func TestDownload(t *testing.T) {
	config := viper.Init("minio")
	picBucket := config.GetString("pic_bucket")

	minio.InitMinIo()
	err := minio.DownloadFile(picBucket, "img.png", "./minioData/receive/")
	if err != nil {
		klog.Info(err.Error())
	}
}

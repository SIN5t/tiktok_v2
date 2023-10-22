package test

import (
	"github.com/SIN5t/tiktok_v2/cmd/video/service"
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"testing"
)

func TestMinio(t *testing.T) {
	config := viper.Init("minio")
	//picBucket := config.GetString("pic_bucket")
	videoBucket := config.GetString("video_bucket")

	minio.InitMinIo()
	//err := minio.UploadFile(picBucket, "img.png", "./img.png", config.GetString("contentType.video"))
	err := minio.UploadFile(videoBucket, "654321.mp4", "./minioData/654321.mp4", config.GetString("contentType.video"))
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

func TestFfmpeg(t *testing.T) {
	service.GetJpegFromFfmpeg("./temp/video/654321.mp4", "./temp/pic/654321.jpeg", "654321.jpeg", 5)
}

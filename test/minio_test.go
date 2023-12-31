package test

import (
	"fmt"
	"github.com/SIN5t/tiktok_v2/cmd/video/service"
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/net/context"
	"testing"
	"time"
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

func TestGetUrl(t *testing.T) {
	minio.InitMinIo()

	config := viper.Init("minio")
	//picBucket := config.GetString("pic_bucket")
	videoBucket := config.GetString("video_bucket")
	url, err := minio.PresignedGetObjectUrl(context.Background(), videoBucket, "654321.mp4", time.Hour*24*7)
	if err != nil {
		klog.Error(err.Error())
	}
	fmt.Println(url)

}

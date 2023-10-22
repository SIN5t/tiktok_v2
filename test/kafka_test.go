package test

import (
	"github.com/SIN5t/tiktok_v2/pkg/kafka/video"
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"testing"
)

func TestKafkaVideoUpload(t *testing.T) {
	msg := &video.VideoMsg{
		VideoPath: "./temp/video/654321.mp4",
		VideoName: "654321.mp4",
		AuthorId:  646156413,
		Title:     "元神启动2",
	}
	video.ProduceMsg(video.NewProducer(), msg)
}
func TestAlone(t *testing.T) {
	minio.InitMinIo()
	consumer := video.NewConsumer()
	video.ConsumePubActMsg(consumer)
}

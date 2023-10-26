package test

import (
	db "github.com/SIN5t/tiktok_v2/cmd/video/dal"
	"github.com/SIN5t/tiktok_v2/pkg/kafka/video"
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"testing"
)

func TestKafkaVideoUpload(t *testing.T) {
	msg := &video.VideoMsg{
		VideoPath: "./temp/video/654321.mp4",
		VideoName: "654321.mp4",
		AuthorId:  646156416,
		Title:     "元神启动3",
	}
	video.ProduceMsg(video.NewProducer(), msg)
}
func TestAlone(t *testing.T) {
	minio.InitMinIo()
	db.InitMysqlDB()
	db.InitRdb()
	consumer := video.NewConsumer()
	video.ConsumePubActMsg(consumer)
}

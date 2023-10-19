package test

import (
	"github.com/SIN5t/tiktok_v2/pkg/kafka/video"
	"testing"
)

func TestKafkaVideoUpload(t *testing.T) {
	msg := &video.VideoMsg{
		VideoPath: "/temp/",
		VideoName: "4asdcsadc.mp4",
		AuthorId:  646156413,
		Title:     "元神启动2",
	}
	video.ProduceMsg(video.NewProducer(), msg)
}
func TestAlone(t *testing.T) {
	consumer := video.NewConsumer()
	video.ConsumePubActMsg(consumer)
}

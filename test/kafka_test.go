package test

import (
	"github.com/SIN5t/tiktok_v2/pkg/kafka"
	"testing"
)

func TestKafkaVideoUpload(t *testing.T) {
	/*producer := kafka.NewProducer()
	video := &db.Video{
		BaseModel: db.BaseModel{
			ID:         0,
			CreateTime: time.Time{},
			UpdateAt:   time.Time{},
			DeleteAt:   gorm.DeletedAt{},
			IsDeleted:  false,
		},
		AuthorID:   0,
		PlayUrl:    "zxc",
		CoverUrl:   "zxc",
		FavCount:   21,
		ComCount:   12,
		IsFavorite: false,
		Data:       nil,
		Title:      "titititititi",
	}
	kafka.ProduceMsg(producer, video)

	consumer := kafka.NewConsumer()
	kafka.ConsumePubActMsg(consumer)*/
}
func TestAlone(t *testing.T) {
	consumer := kafka.NewConsumer()
	kafka.ConsumePubActMsg(consumer)
}

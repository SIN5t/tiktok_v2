package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

func NewProducer() sarama.SyncProducer {
	kafkaViper := viper.Init("kafka")
	endPoint := fmt.Sprintf("%s:%d", kafkaViper.Get("kafka.host"), kafkaViper.Get("kafka.port"))

	config := sarama.NewConfig()
	producer, err := sarama.NewSyncProducer([]string{endPoint}, config)
	if err != nil {
		klog.Fatal(err.Error())
		return nil
	}
	defer producer.Close() //TODO 优雅关闭，何时关闭等
	return producer
}

func SendMsgToBroker(producer sarama.SyncProducer, msg []byte) {
	sendMsg := &sarama.ProducerMessage{
		Topic:     "",
		Key:       nil,
		Value:     nil,
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	partition, offset, err := producer.SendMessage(sendMsg)
	if err != nil {
		klog.Fatal(err.Error())
	}
	klog.Info("Message sent to partition %d at offset %d\n", partition, offset)
}

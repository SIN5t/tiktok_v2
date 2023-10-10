package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

func NewProducer() sarama.SyncProducer {
	kafkaViper := viper.Init("kafka")
	endPoint := fmt.Sprintf("%s:%d", kafkaViper.Get("kafka.host"), kafkaViper.Get("kafka.port"))
	brokers := []string{endPoint}

	// 创建一个客户端
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = config.KafkaVideoClientId
	client, err := sarama.NewClient(brokers, kafkaConfig)
	if err != nil {
		klog.Fatal(err.Error())
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		klog.Fatal(err.Error())
		return nil
	}
	defer producer.Close() //TODO 优雅关闭，何时关闭等
	return producer
}

func SendMsgToBroker(producer sarama.SyncProducer, msg interface{}) {
	msgMarshal, _ := json.Marshal(msg)
	sendMsg := &sarama.ProducerMessage{
		Topic:     config.KafkaVideoTopic,
		Key:       nil,
		Value:     sarama.ByteEncoder(msgMarshal),
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

func NewConsumer() {

}

func ConsumePublishMsg() {

}

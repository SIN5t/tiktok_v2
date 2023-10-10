package video

import (
	"fmt"
	"github.com/IBM/sarama"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
)

import (
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

var (
	kafkaViper = viper.Init("kafka")
	endPoint   = fmt.Sprintf("%s:%d", kafkaViper.Get("kafka.host"), kafkaViper.Get("kafka.port"))
	brokers    = []string{endPoint}
)

func NewProducer() sarama.SyncProducer {

	// 创建一个客户端
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = config.KafkaVideoProducer
	kafkaConfig.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokers, kafkaConfig)
	if err != nil {
		klog.Fatal(err.Error())
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		klog.Fatal(err.Error())
		return nil
	}
	//defer producer.Close() //TODO 优雅关闭，何时关闭等
	return producer
}

func SendMsg(producer sarama.SyncProducer, msg interface{}) {
	msgMarshal, _ := json.Marshal(msg)
	sendMsg := &sarama.ProducerMessage{
		Topic:     config.KafkaVideoTopic,
		Key:       nil,
		Value:     sarama.ByteEncoder(msgMarshal),
		Headers:   nil,
		Metadata:  nil,
		Timestamp: time.Time{},
	}

	partition, offset, err := producer.SendMessage(sendMsg)
	if err != nil {
		klog.Fatal(err.Error())
	}
	klog.Info("Message sent to partition %d at offset %d\n", partition, offset)
}

func NewConsumer() sarama.Consumer {

	//创建客户端
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = config.KafkaVideoConsumer
	kafkaConfig.Producer.Return.Successes = true
	client, err := sarama.NewClient(brokers, kafkaConfig)
	if err != nil {
		klog.Fatal(err.Error())
	}

	//创建消费
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		klog.Fatal(err.Error())
	}
	defer consumer.Close() // TODO 关闭时机
	return consumer
}

func ConsumeMsg(consumer sarama.Consumer) {
	//创建分区消费者
	partitionConsumer, err := consumer.ConsumePartition(config.KafkaVideoTopic, 0, sarama.OffsetNewest)
	if err != nil {
		klog.Fatal(err.Error())
	}
	defer partitionConsumer.Close()
	// 处理消息
	for message := range partitionConsumer.Messages() {
		fmt.Printf("Received message: %s\n", string(message.Value))

	}
}

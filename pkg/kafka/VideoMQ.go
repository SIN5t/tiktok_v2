package kafka

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
	// 等待服务器所有副本都保存成功后的响应
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	//kafkaConfig.ClientID = config.KafkaVideoProducer
	//kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // TODO 待完善
	// 是否等待成功和失败后的响应，只有上面的RequireAcks设置不是NoReponse这里才有用
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
	//defer producer.Close() //TODO 优雅关闭，何时关闭等 : 发送完毕最后一条消息关闭
	return producer
}

func ProduceMsg(producer sarama.SyncProducer, msg interface{}) {
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
	defer producer.Close() // TODO 实际上可以：当发送完最后一条消息后，我们可以调用生产者的 Close() 方法进行关闭

}

func NewConsumer() sarama.Consumer {

	//创建客户端
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	//kafkaConfig.ClientID = config.KafkaVideoConsumer
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
	//defer consumer.Close() // TODO 关闭时机：在程序退出前，或者在接收到一个特定的信号后（例如 SIGINT 或 SIGTERM），调用消费者的 Close() 方法进行关闭
	return consumer
}

func ConsumePubActMsg(consumer sarama.Consumer) {

	// 拿到当前的Topic的所有partition
	partitions, err := consumer.Partitions(config.KafkaVideoTopic)
	if err != nil {
		klog.Fatal(err.Error())
	}

	endChan := make(chan bool)
	for _, partition := range partitions {

		//创建当前分区的消费者
		partitionConsumer, err := consumer.ConsumePartition(config.KafkaVideoTopic, partition, sarama.OffsetOldest)
		if err != nil {
			klog.Fatalf("Failed to create consumer for partition %d: %s", partition, err)
		}

		// 处理消息
		go func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()
			for msg := range pc.Messages() {
				fmt.Printf("Received message: %s\n", string(msg.Value))
			}
			endChan <- true
		}(partitionConsumer)
	}
	<-endChan              // 小心死锁，上面endChan <- true 这里才会执行
	defer consumer.Close() // TODO 实际需改
}

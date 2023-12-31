package video

import (
	"fmt"
	"github.com/IBM/sarama"
	db "github.com/SIN5t/tiktok_v2/cmd/video/dal"
	"github.com/SIN5t/tiktok_v2/cmd/video/service"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/pkg/minio"
	"github.com/SIN5t/tiktok_v2/pkg/snowflakes"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"log"
	"strings"
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

type VideoMsg struct {
	VideoPath string
	VideoName string
	AuthorId  int64
	Title     string
}

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
		partitionConsumer, err := consumer.ConsumePartition(config.KafkaVideoTopic, partition, sarama.OffsetNewest)
		if err != nil {
			klog.Fatalf("Failed to create consumer for partition %d: %s", partition, err)
		}

		// 处理消息
		go func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()
			for msg := range pc.Messages() {
				/*pc.Messages() 是一个阻塞调用
				//当消费者从 pc.Messages() 消费到一条消息时，会执行相应的处理逻辑，然后再次进入循环等待下一条消息。
				//只有当分区消费者 pc 被关闭（通过 defer pc.AsyncClose()）或者出现错误时，才会退出这个循*/

				fmt.Printf("Received message: %s\n, start to processing msg...\n", string(msg.Value))
				videoMsg := &VideoMsg{}
				err := json.Unmarshal(msg.Value, videoMsg)
				if err != nil {
					log.Fatal(err.Error())
				}
				//每个消息开goroutine并行执行
				go func() {

					// 将视频封面等上传OSS
					if UploadToOSS(videoMsg) {
						return
					}

					// 确保视频上传完毕后，再开并行的goroutine，执行可并行的任务
					videoInit := viper.Init("video")
					videoNode := videoInit.GetInt64("snowflake.node")
					dbVideo := &db.Video{
						BaseModel: db.BaseModel{
							ID:         snowflakes.GenerateSnowFlakeId(videoNode),
							CreateTime: time.Now(),
							IsDeleted:  false,
						},
						AuthorID:  videoMsg.AuthorId,
						VideoName: videoMsg.VideoName,
						VideoPath: videoMsg.VideoPath,
						Title:     videoMsg.Title,
					}

					// 上传redis
					go func() {
						// createTime,url需要添加
						db.RedisAddVideo(dbVideo)
						return
					}()
					// 上传mysql
					go func() {
						err = db.SaveVideoToMysql(dbVideo)
						if err != nil {
							log.Fatal(err.Error())
						}
						return
					}()
				}()
			}
			// 设置超时时间自动关闭？ 还是应该一直开着
			endChan <- true
		}(partitionConsumer)
	}

	<-endChan              // 小心死锁，上面 endChan <- true 这里才会执行
	defer consumer.Close() // TODO 实际需改
}

func UploadToOSS(videoMsg *VideoMsg) bool {
	// 上传OSS
	minioViper := viper.Init("minio")
	videoViper := viper.Init("video")
	contentType := minioViper.GetString("contentType.video")

	picName := strings.TrimSuffix(videoMsg.VideoName, "mp4") + "jpeg"
	tempJpegFrameLoc := videoViper.GetString("location.pic") + picName

	err := minio.UploadFile(minioViper.GetString("video_bucket"), videoMsg.VideoName, videoMsg.VideoPath, contentType)
	if err != nil {
		klog.Errorf("视频上传至OSS失败: %s", err.Error())
		return true
	}

	// 截帧，上传图片
	service.GetJpegFromFfmpeg(videoMsg.VideoPath, tempJpegFrameLoc, picName, 10)
	err = minio.UploadFile(minioViper.GetString("pic_bucket"), picName, tempJpegFrameLoc, contentType)
	if err != nil {
		klog.Fatalf("视频帧图片上传到OSS失败：", err.Error())
		return true
	}

	// 所有操作完成之后，删除临时文件
	defer func() {
		// TODO 测试阶段，不删除测试视频
		//err := os.Remove(videoMsg.VideoPath)
		//err = os.Remove(tempJpegFrameLoc)
		if err != nil {
			klog.Errorf("删除文件失败: %s", err.Error())
			return
		} // 文件路径包含文件名
	}()
	return false
}

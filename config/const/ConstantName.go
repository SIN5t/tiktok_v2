package config

const (
	TempVideoLocation  = "./temp/"
	FileAuth           = 0644
	KafkaVideoTopic    = "Topic:publishVideo"
	KafkaVideoProducer = "KafkaClient:VideoMsgProducerId"
	KafkaVideoConsumer = "KafkaClient:VideoMsgConsumerId"
)

const (
	Success = iota
	FailInvalidatePara
	FailInvalidateToken
	FailResponse
)

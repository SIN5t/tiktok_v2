package config

const (
	TempVideoLocation  = "./temp/"
	FileAuth           = 0644
	KafkaVideoTopic    = "Topic:publishVideo"
	KafkaVideoClientId = "KafkaClient:VideoId"
)

const (
	Success = iota
	FailInvalidatePara
	FailInvalidateToken
	FailResponse
)

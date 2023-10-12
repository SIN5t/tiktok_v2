package config

const (
	TempVideoLocation  = "./temp/"
	FileAuth           = 0644
	KafkaVideoTopic    = "publishVideo"
	KafkaVideoProducer = "123456789"
	KafkaVideoConsumer = "12345678910"
)

const (
	Success = iota
	FailInvalidatePara
	FailInvalidateToken
	FailResponse
)

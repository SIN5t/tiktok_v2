package config

const (
	FileAuth           = 0644
	KafkaVideoTopic    = "publishVideo"
	KafkaVideoProducer = "123456789"
	KafkaVideoConsumer = "12345678910"
	JwtSigningKey      = "signing-key"
	JwtIdentityKey     = "userId"
	RedisZsetVideoKey  = "redisVideoKey"
)

const (
	Success = iota
	FailInvalidatePara
	FailInvalidateToken
	FailResponse
)

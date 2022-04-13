package config

type Config struct {
	S3AccessKey                            string
	S3PrivateKey                           string
	KafkaImageClassificationForwardTopic   string
	KafkaImageClassificationResultTopic    string
	KafkaImageClassificationForwardGroupId string
	KafkaBootstrapServers                  string
	KafkaUsername                          string
	KafkaPassword                          string
	KafkaSecurityProtocol                  string
	KafkaSaslMechanism                     string
}

func NewConfig() *Config {
	return &Config{
		S3AccessKey:                            "dummy",
		S3PrivateKey:                           "dummy",
		KafkaImageClassificationForwardTopic:   "dummy",
		KafkaImageClassificationResultTopic:    "dummy",
		KafkaImageClassificationForwardGroupId: "dummy",
		KafkaBootstrapServers:                  "dummy",
		KafkaUsername:                          "dummy",
		KafkaPassword:                          "dummy",
		KafkaSecurityProtocol:                  "dummy",
		KafkaSaslMechanism:                     "dummy",
	}
}

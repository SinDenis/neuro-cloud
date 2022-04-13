package config

import "github.com/confluentinc/confluent-kafka-go/kafka"

func NewImageClassificationProducer(config *Config) *kafka.Producer {
	kafkaConfigMap := kafka.ConfigMap{
		"metadata.broker.list": config.KafkaBootstrapServers,
		"sasl.username":        config.KafkaUsername,
		"sasl.password":        config.KafkaPassword,
		"security.protocol":    config.KafkaSecurityProtocol,
		"sasl.mechanism":       config.KafkaSaslMechanism,
	}

	producer, err := kafka.NewProducer(&kafkaConfigMap)
	if err != nil {
		panic(err)
	}

	return producer
}

func NewImgClassifyResultConsumer(config *Config) *kafka.Consumer {
	kafkaConfigMap := kafka.ConfigMap{
		"group.id":                        config.KafkaImageClassificationResultGroupId,
		"auto.offset.reset":               "smallest",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"metadata.broker.list":            config.KafkaBootstrapServers,
		"sasl.username":                   config.KafkaUsername,
		"sasl.password":                   config.KafkaPassword,
		"security.protocol":               config.KafkaSecurityProtocol,
		"sasl.mechanism":                  config.KafkaSaslMechanism,
	}

	consumer, err := kafka.NewConsumer(&kafkaConfigMap)
	if err != nil {
		panic(err)
	}

	return consumer
}

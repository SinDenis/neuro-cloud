package producer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"nc-file-processing/internal/config"
	"nc-file-processing/internal/model"
	"strconv"
)

type ImageClassifyResultProducer struct {
	logger        logrus.FieldLogger
	kafkaProducer *kafka.Producer
	topic         string
}

func NewImageProducer(config *config.Config, producer *kafka.Producer) *ImageClassifyResultProducer {
	return &ImageClassifyResultProducer{
		logger:        logrus.New(),
		topic:         config.KafkaImageClassificationResultTopic,
		kafkaProducer: producer,
	}
}

func (p *ImageClassifyResultProducer) SendImageClassifyResultEvent(result model.KafkaClassifyImgResult) {
	message, err := convertImageToKafkaMessage(p.topic, result)
	if err != nil {
		p.logger.Error("Failed convert image to kafka message", err)
	}

	deliveryChan := make(chan kafka.Event)
	err = p.kafkaProducer.Produce(message, deliveryChan)
	if err != nil {
		p.logger.Error("Failed produce image classification result", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		p.logger.Error("Failed delivery message to neuro cloud portal", m.TopicPartition.Error)
		return
	}

	p.logger.Infof("Delivered message to topic %s [%d] at offset %v\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
}

func convertImageToKafkaMessage(topic string, result model.KafkaClassifyImgResult) (*kafka.Message, error) {
	key := []byte(strconv.FormatInt(result.ImageId, 10))
	value, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	topicPartition := kafka.TopicPartition{
		Topic:     &topic,
		Partition: kafka.PartitionAny,
	}
	return &kafka.Message{
		TopicPartition: topicPartition,
		Key:            key,
		Value:          value,
	}, nil
}

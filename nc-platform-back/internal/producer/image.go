package producer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"nc-platform-back/internal/config"
	"nc-platform-back/internal/domain"
	"strconv"
)

type ImageProducer struct {
	logger        logrus.FieldLogger
	kafkaProducer *kafka.Producer
	topic         string
}

func NewImageProducer(config *config.Config, producer *kafka.Producer) *ImageProducer {
	return &ImageProducer{
		logger:        logrus.New(),
		topic:         config.KafkaImageClassificationForwardTopic,
		kafkaProducer: producer,
	}
}

func (p *ImageProducer) SendImageClassificationEvent(image domain.Image) {
	message, err := convertImageToKafkaMessage(p.topic, image)
	if err != nil {
		p.logger.Error("Failed convert image to kafka message", err)
	}

	deliveryChan := make(chan kafka.Event)
	err = p.kafkaProducer.Produce(message, deliveryChan)
	if err != nil {
		p.logger.Error("Failed produce image classification message", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		p.logger.Error("Failed delivery message to image event processing service", m.TopicPartition.Error)
		return
	}

	p.logger.Infof("Delivered message to topic %s [%d] at offset %v\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
}

func convertImageToKafkaMessage(topic string, image domain.Image) (*kafka.Message, error) {
	key := []byte(strconv.FormatInt(image.Id, 10))
	imageMsg := imageClassificationMsg{
		ImageId:   image.Id,
		ImageName: image.Name,
		S3ImgLink: image.S3Link,
	}
	value, err := json.Marshal(imageMsg)
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

type imageClassificationMsg struct {
	ImageId   int64  `json:"image_id"`
	ImageName string `json:"image_name"`
	S3ImgLink string `json:"s3_img_link"`
}

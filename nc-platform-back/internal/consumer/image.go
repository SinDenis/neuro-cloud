package consumer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"nc-platform-back/internal/config"
	"nc-platform-back/internal/domain"
	"nc-platform-back/internal/service"
	"os"
	"os/signal"
	"syscall"
)

type ImageClassResultConsumer struct {
	logger        logrus.FieldLogger
	kafkaConsumer *kafka.Consumer
	topic         string
	imgService    *service.ImageService
}

func NewImageConsumer(
	imageClassificationConsumer *kafka.Consumer,
	config *config.Config,
	imgService *service.ImageService,
) *ImageClassResultConsumer {
	return &ImageClassResultConsumer{
		logger:        logrus.New(),
		kafkaConsumer: imageClassificationConsumer,
		topic:         config.KafkaImageClassificationResultTopic,
		imgService:    imgService,
	}
}

func (c *ImageClassResultConsumer) Consume() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	if err := c.kafkaConsumer.Subscribe(c.topic, nil); err != nil {
		c.logger.Fatal(err)
	}

	defer func(kafkaConsumer *kafka.Consumer) {
		err := kafkaConsumer.Close()
		if err != nil {
			c.logger.Error("Failed close kafka consumer", err)
		}
		c.logger.Info("Closing consumer")
	}(c.kafkaConsumer)

	for isKafkaConsumeRun := true; isKafkaConsumeRun; {
		select {
		case sig := <-sigchan:
			c.logger.Info("Caught signal %v: terminating\n", sig)
			isKafkaConsumeRun = false
		case ev := <-c.kafkaConsumer.Events():
			c.logger.Info("New event " + ev.String())
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				if err := c.kafkaConsumer.Assign(e.Partitions); err != nil {
					c.logger.Error("Failed assigned partition", err)
				}
			case kafka.RevokedPartitions:
				if err := c.kafkaConsumer.Unassign(); err != nil {
					c.logger.Error("Failed unassign from partition", err)
				}
			case *kafka.Message:
				c.logger.Infof("%% Message on %s: %s", e.TopicPartition, string(e.Value))
				commit, err := c.kafkaConsumer.Commit()
				if err != nil {
					c.logger.Error("Error commit from kafka consumer", err)
				}

				classifyResult, err := convertKafkaMsg(e.Value)
				if err != nil {
					c.logger.Error("Failed convert kafka msg", err)
				}

				c.imgService.UpdateImage(classifyResult)
				c.logger.Infof("Successful commit %s", commit)
			case kafka.PartitionEOF:
				c.logger.Infof("%% Reached %v", e)
			case kafka.Error:
				c.logger.Errorf("%% Error: %v", e)
				isKafkaConsumeRun = false
			}
		}
	}
}

func convertKafkaMsg(msg []byte) (domain.ClassifyImgResult, error) {
	var classifyResult domain.ClassifyImgResult
	err := json.Unmarshal(msg, &classifyResult)
	if err != nil {
		return domain.ClassifyImgResult{}, err
	}
	return classifyResult, nil
}

package consumer

import (
	"bytes"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"nc-file-processing/internal/config"
	"nc-file-processing/internal/model"
	"nc-file-processing/internal/producer"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ImageConsumer struct {
	logger                 logrus.FieldLogger
	kafkaConsumer          *kafka.Consumer
	topic                  string
	classifyResultProducer *producer.ImageClassifyResultProducer
}

func NewImageConsumer(
	imageClassificationConsumer *kafka.Consumer,
	config *config.Config,
	producer *producer.ImageClassifyResultProducer,
) *ImageConsumer {
	return &ImageConsumer{
		logger:                 logrus.New(),
		kafkaConsumer:          imageClassificationConsumer,
		topic:                  config.KafkaImageClassificationForwardTopic,
		classifyResultProducer: producer,
	}
}

func (c *ImageConsumer) Consume() {
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
				err := c.classifyImg(e.Value)
				if err != nil {
					c.logger.Error(err)
				}
				commit, err := c.kafkaConsumer.Commit()
				if err != nil {
					c.logger.Error("Error commit from kafka consumer", err)
				}
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

func (c *ImageConsumer) classifyImg(kafkaMsg []byte) error {
	classifyImg := model.ClassifyImageKafkaMessage{}
	err := json.Unmarshal(kafkaMsg, &classifyImg)
	if err != nil {
		c.logger.Error("Failed unmarshal kafka message", err)
		return err
	}
	c.logger.Info("Kafka message ", classifyImg)
	response, err := http.Get(classifyImg.S3ImgLink)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	c.logger.Info("Get img from s3 ", response)

	defer response.Body.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("img", classifyImg.ImageName)
	_, err = io.Copy(part, response.Body)
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		c.logger.Error("Failed close writer ", err)
		return err
	}

	r, _ := http.NewRequest("POST", "https://img-neuro-sindenis.cloud.okteto.net/predict", body)
	r.Header.Add("Content-type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		c.logger.Error("Failed send request for img classification in neural network ", err)
		return err
	}
	c.logger.Info("Successful get response from neural network ", resp)
	var classifyImgResult model.NeuralNetClassifyImgResult
	err = json.NewDecoder(resp.Body).Decode(&classifyImgResult)
	if err != nil {
		c.logger.Error("Failed decode classify img result ", err)
		return err
	}

	c.logger.Info("Classify img result ", classifyImgResult)
	c.classifyResultProducer.SendImageClassifyResultEvent(model.KafkaClassifyImgResult{
		ImageId:        classifyImg.ImageId,
		ImageClassName: classifyImgResult.ClassName,
	})
	return nil
}

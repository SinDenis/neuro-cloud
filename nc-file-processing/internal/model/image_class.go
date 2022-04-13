package model

type ClassifyImageKafkaMessage struct {
	ImageId   int64  `json:"image_id"`
	ImageName string `json:"image_name"`
	S3ImgLink string `json:"s3_img_link"`
}

type NeuralNetClassifyImgResult struct {
	ClassId   string `json:"class_id"`
	ClassName string `json:"class_name"`
}

type KafkaClassifyImgResult struct {
	ImageId        int64  `json:"image_id"`
	ImageClassName string `json:"image_class_name"`
	Category       string `json:"category"`
}

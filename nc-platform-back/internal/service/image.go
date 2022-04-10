package service

import (
	"context"
	"demo-rest/internal/config"
	"demo-rest/internal/domain"
	"demo-rest/internal/dto"
	"demo-rest/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"time"
)

type ImageService struct {
	logger          logrus.FieldLogger
	s3Uploader      *s3manager.Uploader
	imageRepository *repository.ImageRepository
}

func NewImageService(config *config.Config, imageRepository *repository.ImageRepository) *ImageService {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-2"),
		Credentials: credentials.NewStaticCredentials(config.S3AccessKey, config.S3PrivateKey, ""),
	})
	if err != nil {
		panic(err)
	}

	return &ImageService{
		logger:          logrus.New(),
		s3Uploader:      s3manager.NewUploader(sess),
		imageRepository: imageRepository,
	}
}

func (s *ImageService) GetUserImages(ctx context.Context, pagingParam dto.PagingParam) ([]domain.Image, error) {
	userId := int64(ctx.Value("userId").(float64))
	s.logger.Info(userId)
	return s.imageRepository.GetUserImages(userId, pagingParam)
}

func (s *ImageService) Upload(context context.Context, file multipart.File, header *multipart.FileHeader) error {
	response, err := s.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("dsin-neuro-storage"),
		Key:    aws.String("images/" + header.Filename),
		Body:   file,
	})
	if err != nil {
		return err
	}

	img := domain.Image{
		Name:         header.Filename,
		Size:         header.Size,
		S3Link:       response.Location,
		UserId:       int64(context.Value("userId").(float64)),
		DateUploaded: time.Now().UTC(),
	}
	err = s.imageRepository.SaveImg(img)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"nc-platform-back/internal/domain"
	"nc-platform-back/internal/dto"
)

const (
	selectNextId   = "select nextval('image_seq')"
	insertImgQuery = `insert into image(id, name, size, description, s3_link, dateUploaded, user_id)
						values ($1, $2, $3, $4, $5, $6, $7)`
	selectUserImagesQuery = `select id, name, size, s3_link, dateuploaded, label
							from image where user_id = $1 order by dateuploaded desc limit $2 offset $3`
	updateImgLabelById = "update image set label = $1 where id = $2"
)

type ImageRepository struct {
	logger logrus.FieldLogger
	pool   *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{logger: logrus.New(), pool: pool}
}

func (r *ImageRepository) SaveImg(image domain.Image) (domain.Image, error) {
	var description interface{} = nil
	if len(image.Description) > 0 {
		description = image.Description
	}
	id, err := r.getNextId()
	image.Id = id
	if err != nil {
		return domain.Image{}, err
	}
	query, err := r.pool.Query(
		context.Background(),
		insertImgQuery,
		image.Id, image.Name, image.Size, description, image.S3Link, image.DateUploaded, image.UserId)
	defer query.Close()
	if err != nil {
		return domain.Image{}, err
	}
	return image, nil
}

func (r *ImageRepository) GetUserImages(userId int64, param dto.PagingParam) ([]domain.Image, error) {
	offset := param.PageSize * (param.Page - 1)
	rows, err := r.pool.Query(context.Background(), selectUserImagesQuery, userId, param.PageSize, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var images []domain.Image
	var image domain.Image
	for rows.Next() {
		var imageLabel sql.NullString
		err := rows.Scan(&image.Id, &image.Name, &image.Size, &image.S3Link, &image.DateUploaded, &imageLabel)
		if imageLabel.Valid {
			image.Label = imageLabel.String
		}
		if err != nil {
			fmt.Println(err)
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *ImageRepository) getNextId() (int64, error) {
	var id int64
	rows, _ := r.pool.Query(context.Background(), selectNextId)
	defer rows.Close()
	rows.Next()
	rows.Scan(&id)
	return id, nil
}

func (r *ImageRepository) UpdateImgLabel(imgId int64, label string) {
	rows, err := r.pool.Query(context.Background(), updateImgLabelById, label, imgId)
	if err != nil {
		r.logger.Errorf("Failed update label for img = %d", imgId, err)
	}
	rows.Close()
}

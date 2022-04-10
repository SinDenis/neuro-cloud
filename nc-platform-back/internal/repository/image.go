package repository

import (
	"context"
	"demo-rest/internal/domain"
	"demo-rest/internal/dto"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertImgQuery = `insert into image(id, name, size, description, s3_link, dateUploaded, user_id)
						values (nextval('image_seq'), $1, $2, $3, $4, $5, $6)`
	selectUserImagesQuery = `select id, name, size, s3_link, dateuploaded
							from image where user_id = $1 order by dateuploaded desc limit $2 offset $3`
)

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) SaveImg(image domain.Image) error {
	var description interface{} = nil
	if len(image.Description) > 0 {
		description = image.Description
	}
	query, err := r.pool.Query(
		context.Background(),
		insertImgQuery,
		image.Name, image.Size, description, image.S3Link, image.DateUploaded, image.UserId)
	defer query.Close()
	if err != nil {
		return err
	}
	return nil
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
		err := rows.Scan(&image.Id, &image.Name, &image.Size, &image.S3Link, &image.DateUploaded)
		if err != nil {
			fmt.Println(err)
		}
		images = append(images, image)
	}
	return images, nil
}

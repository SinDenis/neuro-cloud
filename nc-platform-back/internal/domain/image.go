package domain

import "time"

type Image struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Size         int64     `json:"size"`
	S3Link       string    `json:"s3Link"`
	Label        string    `json:"label"`
	Category     string    `json:"category"`
	DateUploaded time.Time `json:"dateUploaded"`
	UserId       int64     `json:"userId"`
}

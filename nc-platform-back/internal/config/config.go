package config

import (
	"demo-rest/pkg/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Dsn          string
	JwtKey       string
	JwtAlgo      string
	S3AccessKey  string
	S3PrivateKey string
}

func NewConfig() *Config {
	return &Config{
		Dsn:          "dummy",
		JwtKey:       "dummy",
		JwtAlgo:      "dummy",
		S3AccessKey:  "dummy",
		S3PrivateKey: "dummy",
	}
}

func NewPgPool(config *Config) *pgxpool.Pool {
	pool, err := postgres.NewPool(config.Dsn)
	if err != nil {
		panic(err)
	}
	return pool
}

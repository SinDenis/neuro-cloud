package config

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"nc-platform-back/pkg/postgres"
)

type Config struct {
	Dsn                                   string
	JwtKey                                string
	JwtAlgo                               string
	S3AccessKey                           string
	S3PrivateKey                          string
	KafkaImageClassificationForwardTopic  string
	KafkaImageClassificationResultTopic   string
	KafkaImageClassificationResultGroupId string
	KafkaBootstrapServers                 string
	KafkaUsername                         string
	KafkaPassword                         string
	KafkaSecurityProtocol                 string
	KafkaSaslMechanism                    string
}

func NewConfig() *Config {
	return &Config{
		Dsn:                                   "dummy",
		JwtKey:                                "dummy",
		JwtAlgo:                               "dummy",
		S3AccessKey:                           "dummy",
		S3PrivateKey:                          "dummy",
		KafkaImageClassificationForwardTopic:  "dummy",
		KafkaImageClassificationResultTopic:   "dummy",
		KafkaImageClassificationResultGroupId: "dummy",
		KafkaBootstrapServers:                 "dummy",
		KafkaUsername:                         "dummy",
		KafkaPassword:                         "dummy",
		KafkaSecurityProtocol:                 "dummy",
		KafkaSaslMechanism:                    "dummy",
	}
}

func NewPgPool(config *Config) *pgxpool.Pool {
	pool, err := postgres.NewPool(config.Dsn)
	if err != nil {
		panic(err)
	}
	return pool
}

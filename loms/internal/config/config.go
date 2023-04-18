package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type CancelReservationDueTimeoutJob struct {
	Cron                  string        `yaml:"cron"`
	OrderToBePayedTimeout time.Duration `yaml:"orderToBePayedTimeout"`
}

type ReadOutBoxSendJob struct {
	Cron            string `yaml:"cron"`
	BatchSizeToRead int    `yaml:"batchSizeToRead"`
	TopicToSend     string `yaml:"topicToSend"`
}

type CronJobs struct {
	ReadOutBoxSendJob              ReadOutBoxSendJob              `yaml:"readOutboxSendJob"`
	CancelReservationDueTimeoutJob CancelReservationDueTimeoutJob `yaml:"cancelReservationDueTimeoutJob"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	DBName   string `yaml:"dbName"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	SslMode  string `yaml:"sslMode"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers"`
}

type Config struct {
	ServiceName string   `yaml:"serviceName"`
	LoggerLevel string   `yaml:"loggerLevel"`
	GrpcPort    int      `yaml:"grpcPort"`
	MetricsPort int      `yaml:"metricsPort"`
	DbConfig    DbConfig `yaml:"dbConfig"`
	CronJobs    CronJobs `yaml:"cronJobs"`
	Kafka       Kafka    `yaml:"kafka"`
}

var Data Config

func New() error {
	rawYAML, err := os.ReadFile("config.yaml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &Data)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}

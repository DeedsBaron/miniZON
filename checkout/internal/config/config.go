package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName string `yaml:"serviceName"`
	MetricsPort int    `yaml:"metricsPort"`
	GrpcPort    int    `yaml:"grpcPort"`
	Token       string `yaml:"token"`
	Services    struct {
		Loms string `yaml:"loms"`
		Ps   string `yaml:"ps"`
	} `yaml:"services"`
	DbConfig    DbConfig `yaml:"dbConfig"`
	Workers     Workers  `yaml:"workers"`
	LoggerLevel string   `yaml:"loggerLevel"`
	Cache       Cache    `yaml:"cache"`
}

type Cache struct {
	Ttl         time.Duration `yaml:"ttl"`
	Buckets     int           `yaml:"buckets"`
	LruCapacity int           `yaml:"lruCapacity"`
}
type DbConfig struct {
	Host     string `yaml:"host"`
	DBName   string `yaml:"dbName"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	SslMode  string `yaml:"sslMode"`
}

type Workers struct {
	Amount int `yaml:"amount"`
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

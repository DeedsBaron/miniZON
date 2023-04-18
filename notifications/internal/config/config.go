package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Kafka struct {
	Brokers       []string `yaml:"brokers"`
	Topic         string   `yaml:"topic"`
	ConsumerGroup string   `yaml:"consumerGroup"`
	Workers       int      `yaml:"workers"`
}

type Config struct {
	Kafka       Kafka  `yaml:"kafka"`
	LoggerLevel string `yaml:"loggerLevel"`
	MetricsPort int    `yaml:"metricsPort"`
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

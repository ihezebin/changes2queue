package config

import (
	"encoding/json"
	"os"

	"github.com/ihezebin/soup/config"
	"github.com/ihezebin/soup/logger"
	"github.com/pkg/errors"

	"github.com/ihezebin/changes2queue/component/task"
)

type Config struct {
	ServiceName  string        `json:"service_name" mapstructure:"service_name"`
	Port         uint          `json:"port" mapstructure:"port"`
	Logger       *LoggerConfig `json:"logger" mapstructure:"logger"`
	Mongo2Queues task.Config   `json:"mongo2queues" mapstructure:"mongo2queues"`
	MySQL2Queues task.Config   `json:"mysql2queues" mapstructure:"mysql2queues"`
	Pwd          string        `json:"-" mapstructure:"-"`
}

type LoggerConfig struct {
	Level    logger.Level `json:"level" mapstructure:"level"`
	Filename string       `json:"filename" mapstructure:"filename"`
}

var gConfig *Config = &Config{}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func GetConfig() *Config {
	return gConfig
}

func Load(path string) (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "get pwd error")
	}

	if err = config.NewWithFilePath(path).Load(gConfig); err != nil {
		return nil, errors.Wrap(err, "load config error")
	}

	gConfig.Pwd = pwd

	return gConfig, nil
}

package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config - основная структура конфига приложения которая агрегирует все другие конфиги
type Config struct {
	App App
}

// GetConfig - получает конфиг файла на основе переменных окружения
func GetConfig() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("error while init config: %s", err.Error())
	}

	return config, nil
}

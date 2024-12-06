package config

import (
	"xdb/internal/shared"

	"github.com/gofiber/fiber/v2/log"
)

func LoadConfig(filepath string) (*GlobalConfig, error) {
	config := &GlobalConfig{}

	if err := shared.LoadTomlFile(config, filepath); err != nil {
		log.Error(err)
		return nil, errUnableToLoadConfig
	}

	return config, nil
}

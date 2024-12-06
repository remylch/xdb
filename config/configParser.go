package config

import (
	"xdb/internal/shared"
)

func LoadConfig(filepath string) (*GlobalConfig, error) {
	config := &GlobalConfig{}

	if err := shared.LoadTomlFile(config, filepath); err != nil {
		return nil, errUnableToLoadConfig
	}

	return nil, errUnableToLoadConfig
}

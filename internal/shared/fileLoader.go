package shared

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadTomlFile(configType interface{}, filepath string) error {
	file, err := os.ReadFile(filepath)

	if err != nil {
		return err
	}

	if err := toml.Unmarshal(file, configType); err != nil {
		return errUnableToLoadTomlFile
	}

	return nil
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	defaultTomlConfig := "../config.toml"

	conf, err := LoadConfig(defaultTomlConfig)

	assert.NoError(t, err, "load config should not throw error")

	assert.Equal(t, conf.Server.ApiAddr, ":8080")
	assert.Equal(t, conf.Server.NodeAddr, ":3000")
	assert.Equal(t, conf.Log.LogDir, "./log/")
}

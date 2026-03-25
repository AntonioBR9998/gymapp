package config

import (
	"github.com/AntonioBR9998/go-common/config"
)

// Config represents the service configuration
type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	ServerName string `json:"serverName"`
}

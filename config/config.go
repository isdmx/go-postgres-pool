package config

import (
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Database string
	Host     string
	User     string
	Password string
	Pool     Pool
	TLS      TLSConfig
	Params   map[string]string
}

type Pool struct {
	MinConns              int32
	MaxConns              int32
	MaxConnIdleTime       time.Duration
	MaxConnLifetime       time.Duration
	MaxConnLifetimeJitter time.Duration
	HealthCheckPeriod     time.Duration
}

type TLSConfig struct {
	Enabled    bool
	CACert     string
	ClientCert string
	ClientKey  string
	SkipVerify bool
}

func LoadConfig(path string) (*DBConfig, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg DBConfig
	if err := viper.UnmarshalKey("DB", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

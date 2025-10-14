package config

import (
	"path/filepath"

	"github.com/shlyyy/micro-services/pkg/logger"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`

	Logger logger.LogConfig `mapstructure:"logger"`
}

var Cfg AppConfig

func LoadConfig(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	dir := filepath.Dir(path)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(&Cfg); err != nil {
		return err
	}
	return nil
}

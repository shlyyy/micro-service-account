package config

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/shlyyy/micro-service-account/pkg/db"
	"github.com/shlyyy/micro-service-account/pkg/logger"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Account struct {
		AccountServer struct {
			GrpcServerPort int `mapstructure:"grpc_server_port"`
		} `mapstructure:"account_server"`
		AccountWeb struct {
			Port int `mapstructure:"port"`
		} `mapstructure:"account_web"`
	} `mapstructure:"account"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`

	Database db.DBConfig      `mapstructure:"database"`
	Logger   logger.LogConfig `mapstructure:"logger"`
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

	// 动态监控配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		// 配置文件更新时的回调函数
		logger.Info("配置文件已更新:", in.Name)

		if err := v.ReadInConfig(); err != nil {
			logger.Error("重新读取配置文件失败:", err)
			return
		}
		if err := v.Unmarshal(&Cfg); err != nil {
			logger.Error("重新解析配置文件失败:", err)
			return
		}
		logger.Info("配置文件重新加载成功")
	})
	return nil
}

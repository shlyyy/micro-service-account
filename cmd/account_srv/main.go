package main

import (
	"fmt"

	"github.com/shlyyy/micro-services/pkg/config"
	"github.com/shlyyy/micro-services/pkg/db"
	"github.com/shlyyy/micro-services/pkg/logger"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		panic(fmt.Errorf("加载配置失败: %w", err))
	}

	// 初始化日志
	logger.Init(config.Cfg.Logger)
	logger.Info("服务启动成功")

	// 数据库连接
	db.InitDB(config.Cfg.Database.DSN)
	logger.Debug(db.DB)
}

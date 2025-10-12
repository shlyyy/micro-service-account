package db

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
}

func InitDB(cfg DBConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	newLogger := logger.New(
		zap.NewExample(), // 可替换为 zap.L()
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// 自动迁移（可选）
	// db.AutoMigrate(&model.Account{})

	zap.L().Info("✅ MySQL connected")
	return nil
}

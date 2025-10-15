package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	mylogger "github.com/shlyyy/micro-services/pkg/logger"
)

var DB *gorm.DB

type DBConfig struct {
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Dbname          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	ParseTime       string `mapstructure:"parse_time"`
	Loc             string `mapstructure:"loc"`
	LogPath         string `mapstructure:"log_path"`
	LogLevel        string `mapstructure:"log_level"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// func InitDB(dsn string) error {
func InitDB(cfg *DBConfig) error {
	// 构造 DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)
	mylogger.Info("数据库连接字符串:", dsn)

	// 创建数据库日志文件
	logFile, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		mylogger.Error("Cannot open log file:", cfg.LogPath)
		return err
	}
	// defer logFile.Close()

	// 定义日志级别
	var level logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		level = logger.Silent
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	default:
		level = logger.Info
	}

	// 创建 GORM 日志器
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	gormStdLogger := log.New(multiWriter, "\r\n", log.LstdFlags)

	newLogger := logger.New(
		gormStdLogger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second) // 连接最大生命周期

	DB = db

	return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		mylogger.Error("Database not initialized. Please call InitDB first.")
		return nil
	}
	return DB
}

package db

import (
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

func InitDB(dsn string) error {
	logFilePath := "gorm.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	gormStdLogger := log.New(multiWriter, "\r\n", log.LstdFlags)

	newLogger := logger.New(
		gormStdLogger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		mylogger.Error("MySQL connection failed:", err)
		return err
	}

	DB = db

	mylogger.Info("MySQL connected")
	return nil
}

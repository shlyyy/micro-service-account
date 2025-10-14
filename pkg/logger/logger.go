package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log  *zap.SugaredLogger
	once sync.Once
)

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	Console    bool   `mapstructure:"console"`
	JsonFormat bool   `mapstructure:"json"`
	Color      bool   `mapstructure:"color"`
}

func Init(cfg LogConfig) {
	once.Do(func() {
		level := parseLevel(cfg.Level)
		encoderCfg := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stack",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeCaller: zapcore.ShortCallerEncoder,
		}

		// ---- 文件日志（无颜色）
		fileEncoderCfg := encoderCfg
		fileEncoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder := zapcore.NewConsoleEncoder(fileEncoderCfg)

		// ---- 控制台日志（彩色）
		consoleEncoderCfg := encoderCfg
		consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)

		// ---- 输出目标
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})

		var core zapcore.Core
		if cfg.Console {
			core = zapcore.NewTee(
				zapcore.NewCore(fileEncoder, fileWriter, level),
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)
		} else {
			core = zapcore.NewCore(fileEncoder, fileWriter, level)
		}

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		log = zapLogger.Sugar()
	})
}

func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Debug(args ...interface{}) { log.Debug(args...) }
func Info(args ...interface{})  { log.Info(args...) }
func Warn(args ...interface{})  { log.Warn(args...) }
func Error(args ...interface{}) { log.Error(args...) }

func Debugf(tmpl string, args ...interface{}) { log.Debugf(tmpl, args...) }
func Infof(tmpl string, args ...interface{})  { log.Infof(tmpl, args...) }
func Warnf(tmpl string, args ...interface{})  { log.Warnf(tmpl, args...) }
func Errorf(tmpl string, args ...interface{}) { log.Errorf(tmpl, args...) }

func Sync() { _ = log.Sync() }

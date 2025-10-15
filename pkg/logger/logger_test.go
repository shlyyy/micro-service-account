package logger

import "testing"

func TestLogger(t *testing.T) {
	cfg := LogConfig{
		Level:      "debug",
		Filename:   "test.log",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
		Console:    true,
		JsonFormat: false,
	}

	InitLogger(cfg)
	Info("info message")
	Debug("debug message")
	Error("error message")
	Sync()
}

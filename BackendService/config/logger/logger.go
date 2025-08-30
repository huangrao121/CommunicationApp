package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	LevelVar = new(slog.LevelVar)
)

func InitLogger() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", "error", err)
	}

	logLevel := os.Getenv("LOG_LEVEL")
	logFormat := os.Getenv("LOG_FORMAT")
	logPath := os.Getenv("LOG_PATH")

	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	var writer io.Writer
	if logPath != "" {
		logRotation := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    10,
			MaxBackups: 2,
			MaxAge:     30,
			Compress:   true,
		}
		writer = io.MultiWriter(os.Stdout, logRotation)
	} else {
		writer = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 如果日志级别是 SourceKey，则不输出，隐藏源码位置
			if a.Key == slog.SourceKey {
				return slog.Attr{}
			}
			return a
		},
	}

	// 在开发环境（text格式）中，我们通常希望看到源码位置
	if strings.ToLower(logFormat) == "text" {
		opts.AddSource = true
		opts.ReplaceAttr = nil // 使用默认行为
	}

	var handler slog.Handler
	if strings.ToLower(logFormat) == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}
	// 设置默认的 slog 处理器
	slog.SetDefault(slog.New(handler))
}

func SetLevel(level slog.Level) {
	LevelVar.Set(level)
}

func WithCtx(ctx context.Context) *slog.Logger { return slog.Default().With(slog.Group("ctx")) }

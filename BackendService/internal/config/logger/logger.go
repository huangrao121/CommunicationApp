package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	LevelVar = new(slog.LevelVar)
)

func InitLogger(service, version, env string, level slog.Level) {
	LevelVar.Set(level)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       LevelVar, // 生态级别
		AddSource:   false,    // production时关，排除障碍可以开
		ReplaceAttr: nil,      // 需要映射字段名可以定义
	})
	l := slog.New(h).With(
		"service", service,
		"version", version,
		"env", env,
		// "time", slog.Time(time.Now(), "time"),
	)

	slog.SetDefault(l)
}

func SetLevel(level slog.Level) {
	LevelVar.Set(level)
}

func WithCtx(ctx context.Context) *slog.Logger { return slog.Default().With(slog.Group("ctx")) }

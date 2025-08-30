package main

import (
	"log/slog"

	//"time"

	//"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/config"
	"github.com/huangrao121/CommunicationApp/BackendService/config/database"
	"github.com/huangrao121/CommunicationApp/BackendService/config/logger"

	//"github.com/huangrao121/CommunicationApp/BackendService/internal/user"

	"github.com/huangrao121/CommunicationApp/BackendService/internal/http"
)

func main() {
	// 从config.yaml中加载配置
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}
	// 这个的配置来自env
	logger.InitLogger()

	// 初始化数据库
	database.InitDB(cfg)

	slog.Info("BackendService started")

	// user := user.User{
	// 	ID:        uuid.New().String(),
	// 	Username:  "test",
	// 	Email:     "test@test.com",
	// 	Password:  "test",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }
	// token, err := pkg.GenerateJWKToken(&user)
	// if err != nil {
	// 	slog.Error("failed to generate jwt token", "error", err)
	// }
	// slog.Info("jwt token", "token", token)

	// claims, err := pkg.ParseJWKToken(token)
	// if err != nil {
	// 	slog.Error("failed to parse jwt token", "error", err)
	// }
	// slog.Info("jwt claims", "claims", claims)

	router := http.InitRouter()
	router.Run(":8080")
}

package main

import (
	"log/slog"
	"os"

	//"time"

	//"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/config"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/config/logger"

	//"github.com/huangrao121/CommunicationApp/BackendService/internal/config/pkg"
	//"github.com/huangrao121/CommunicationApp/BackendService/internal/user"

	"github.com/huangrao121/CommunicationApp/BackendService/internal/http"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func main() {
	config.InitDB()
	logger.InitLogger("BackendService", os.Getenv("VERSION"), "dev", 0)
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

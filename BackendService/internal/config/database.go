package config

import (
	"log/slog"
	"sync"

	"os"

	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		dsn := os.Getenv("DB_DSN")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Error("failed to connect database", "error", err)
		}
		result := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		if result.Error != nil {
			slog.Error("failed to create uuid extension", "error", result.Error)
		}
		joinTableErr := db.SetupJoinTable(&types.User{}, "Friendships", &types.Friendship{})
		if joinTableErr != nil {
			slog.Error("failed to setup join table for friendships", "error", joinTableErr)
		}
		migrateErr := db.AutoMigrate(&types.User{}, &types.OauthIdentity{})
		if migrateErr != nil {
			slog.Error("failed to migrate database", "error", migrateErr)
		}
		slog.Info("database migrate successfully")
	})
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}

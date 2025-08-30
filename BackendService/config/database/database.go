package database

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/huangrao121/CommunicationApp/BackendService/config"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB(cfg *config.Config) {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Error("failed to connect database", "error", err)
		}
		result := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
		if result.Error != nil {
			slog.Error("failed to create uuid extension", "error", result.Error)
		}

		migrateErr := db.AutoMigrate(
			&types.Users{},
			&types.OauthIdentities{},
			&types.Groups{},
			&types.Conversations{},
			&types.P2PMessages{},
			&types.GroupMessages{},
			&types.ConversationParticipants{},
			&types.Friends{},
			&types.GroupMembers{},
		)
		if migrateErr != nil {
			slog.Error("failed to migrate database", "error", migrateErr)
		}
		slog.Info("database migrate successfully")
	})
}

func GetDB(cfg *config.Config) *gorm.DB {
	if db == nil {
		InitDB(cfg)
	}
	return db
}

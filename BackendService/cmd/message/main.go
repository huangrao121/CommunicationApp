package main

import (
	"log"

	"github.com/huangrao121/CommunicationApp/BackendService/config"
	"github.com/huangrao121/CommunicationApp/BackendService/config/database"
	"github.com/huangrao121/CommunicationApp/BackendService/config/logger"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/common/kafka"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/message/handler"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/message/service"
)

func main() {
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("failed to load config", "error", err)
		return
	}
	logger.InitLogger()

	// 初始化db
	database.InitDB(cfg)

	// 初始化Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	// message service
	db := database.GetDB(cfg)
	messageService := service.NewMessageService(db, kafkaProducer)
	messageHandler := handler.NewMessageHandler(messageService)
	handlerInit := NewHandlerInit(messageHandler)

	// 初始化gin http
	router := InitializeRouter(handlerInit)
	router.Run(":8081")

}

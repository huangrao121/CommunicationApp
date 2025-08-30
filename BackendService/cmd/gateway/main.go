package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huangrao121/CommunicationApp/BackendService/config"

	"github.com/huangrao121/CommunicationApp/BackendService/internal/common/middleware"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/handler"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/service"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/websocket"
)

func main() {
	// 加载配置, Loadconfig接收路径，使用了相对路径。
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Gateway服务（与Message Service通信）
	gatewayService := service.NewGatewayService("http://message-service:8002")

	// WebSocket Hub
	hub := websocket.NewHub(gatewayService)

	// 启动Hub
	ctx := context.Background()
	go hub.Run(ctx)

	// 处理器
	gatewayHandler := handler.NewGatewayHandler(hub)

	// 设置路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// WebSocket路由
	r.GET("/ws", middleware.AuthMiddleware(), gatewayHandler.HandleWebSocket)

	// API路由
	api := r.Group("/api/v1")
	{
		api.GET("/online-users", gatewayHandler.GetOnlineUsers)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	log.Printf("Gateway service starting on port %d", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

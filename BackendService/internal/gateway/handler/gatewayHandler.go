package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/websocket"
)

type GatewayHandler struct {
	hub *websocket.Hub
}

func NewGatewayHandler(hub *websocket.Hub) *GatewayHandler {
	return &GatewayHandler{
		hub: hub,
	}
}

func (h *GatewayHandler) HandleWebSocket(c *gin.Context) {
	// 从查询参数或头部获取token
	userName := c.GetString("username")
	userID := c.GetString("userID")

	// 处理WebSocket连接
	h.hub.HandleWebSocket(c.Writer, c.Request, uuid.MustParse(userID), userName)
}

func (h *GatewayHandler) GetOnlineUsers(c *gin.Context) {
	// 返回在线用户列表
	onlineUsers := h.hub.GetOnlineUsers()
	c.JSON(http.StatusOK, gin.H{
		"online_users": onlineUsers,
	})
}

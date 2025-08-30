package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients        map[uuid.UUID]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	userMessage    chan UserMessage
	messageService MessageServiceClient // gRPC客户端接口
	mutex          sync.RWMutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uuid.UUID
	username string
}

type UserMessage struct {
	UserID  uuid.UUID `json:"user_id"`
	Type    string    `json:"type"`
	Payload []byte    `json:"payload"`
}

type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// MessageServiceClient 接口，用于与Message Service通信
type MessageServiceClient interface {
	SendP2PMessage(ctx context.Context, req *SendP2PRequest) (*MessageResponse, error)
	SendGroupMessage(ctx context.Context, req *SendGroupRequest) (*MessageResponse, error)
}

type SendP2PRequest struct {
	SenderID    uuid.UUID `json:"sender_id"`
	ReceiverID  uuid.UUID `json:"receiver_id"`
	Content     string    `json:"content"`
	ContentType int       `json:"content_type"`
}

type SendGroupRequest struct {
	SenderID    uuid.UUID `json:"sender_id"`
	GroupID     uuid.UUID `json:"group_id"`
	Content     string    `json:"content"`
	ContentType int       `json:"content_type"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Timestamp int64     `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要更严格的检查
	},
	HandshakeTimeout: 45 * time.Second,
}

func NewHub(messageService MessageServiceClient) *Hub {
	return &Hub{
		clients:        make(map[uuid.UUID]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		userMessage:    make(chan UserMessage),
		messageService: messageService,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
			log.Printf("Client %s (%s) connected", client.username, client.userID)

			// 发送连接成功消息
			client.sendMessage(OutgoingMessage{
				Type:      "connection_established",
				Data:      map[string]interface{}{"user_id": client.userID},
				Timestamp: time.Now().Unix(),
			})

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("Client %s (%s) disconnected", client.username, client.userID)

		case message := <-h.broadcast:
			h.broadcastToAll(message)

		case userMsg := <-h.userMessage:
			h.handleUserMessage(ctx, userMsg)
		}
	}
}

func (h *Hub) handleUserMessage(ctx context.Context, userMsg UserMessage) {
	var incoming IncomingMessage
	if err := json.Unmarshal(userMsg.Payload, &incoming); err != nil {
		log.Printf("Error unmarshaling user message: %v", err)
		return
	}

	switch incoming.Type {
	case "send_p2p_message":
		h.handleP2PMessage(ctx, userMsg.UserID, incoming.Data)
	case "send_group_message":
		h.handleGroupMessage(ctx, userMsg.UserID, incoming.Data)
	case "typing":
		h.handleTyping(userMsg.UserID, incoming.Data)
	case "read_receipt":
		h.handleReadReceipt(userMsg.UserID, incoming.Data)
	default:
		log.Printf("Unknown message type: %s", incoming.Type)
	}
}

func (h *Hub) handleP2PMessage(ctx context.Context, senderID uuid.UUID, data json.RawMessage) {
	var req SendP2PRequest
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("Error unmarshaling P2P message: %v", err)
		return
	}
	req.SenderID = senderID

	// 调用Message Service
	resp, err := h.messageService.SendP2PMessage(ctx, &req)
	if err != nil {
		log.Printf("Error sending P2P message: %v", err)
		h.sendErrorToUser(senderID, "Failed to send message", err)
		return
	}

	// 发送给接收者
	h.SendToUser(req.ReceiverID, OutgoingMessage{
		Type:      "new_p2p_message",
		Data:      resp,
		Timestamp: time.Now().Unix(),
	})

	// 发送确认给发送者
	h.SendToUser(senderID, OutgoingMessage{
		Type:      "message_sent",
		Data:      resp,
		Timestamp: time.Now().Unix(),
	})
}

func (h *Hub) handleGroupMessage(ctx context.Context, senderID uuid.UUID, data json.RawMessage) {
	var req SendGroupRequest
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("Error unmarshaling group message: %v", err)
		return
	}
	req.SenderID = senderID

	// 调用Message Service
	//resp, err := h.messageService.SendGroupMessage(ctx, &req)
	//if err != nil {
	//	log.Printf("Error sending group message: %v", err)
	//	h.sendErrorToUser(senderID, "Failed to send group message", err)
	//	return
	//}

	// 这里可以从Message Service获取群组成员列表，然后推送给所有成员
	// 或者通过Kafka事件来处理群组消息分发
}

func (h *Hub) handleTyping(senderID uuid.UUID, data json.RawMessage) {
	var typingData struct {
		ReceiverID uuid.UUID `json:"receiver_id"`
		IsTyping   bool      `json:"is_typing"`
	}

	if err := json.Unmarshal(data, &typingData); err != nil {
		return
	}

	h.SendToUser(typingData.ReceiverID, OutgoingMessage{
		Type: "typing_indicator",
		Data: map[string]interface{}{
			"user_id":   senderID,
			"is_typing": typingData.IsTyping,
		},
		Timestamp: time.Now().Unix(),
	})
}

func (h *Hub) handleReadReceipt(senderID uuid.UUID, data json.RawMessage) {
	// 处理已读回执逻辑
}

func (h *Hub) SendToUser(userID uuid.UUID, message OutgoingMessage) {
	h.mutex.RLock()
	client, ok := h.clients[userID]
	h.mutex.RUnlock()

	if !ok {
		return
	}

	client.sendMessage(message)
}

func (h *Hub) broadcastToAll(message []byte) {
	h.mutex.RLock()
	for _, client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client.userID)
		}
	}
	h.mutex.RUnlock()
}

func (h *Hub) sendErrorToUser(userID uuid.UUID, message string, err error) {
	errorMsg := OutgoingMessage{
		Type: "error",
		Data: map[string]interface{}{
			"message": message,
			"error":   err.Error(),
		},
		Timestamp: time.Now().Unix(),
	}
	h.SendToUser(userID, errorMsg)
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID uuid.UUID, username string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (h *Hub) GetOnlineUsers() []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	onlineUsers := make([]uuid.UUID, 0, len(h.clients))

	for _, client := range h.clients {
		onlineUsers = append(onlineUsers, client.userID)
	}

	return onlineUsers
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 发送消息到处理器
		c.hub.userMessage <- UserMessage{
			UserID:  c.userID,
			Type:    "user_message",
			Payload: message,
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendMessage(message OutgoingMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case c.send <- data:
	default:
		c.hub.mutex.Lock()
		close(c.send)
		delete(c.hub.clients, c.userID)
		c.hub.mutex.Unlock()
	}
}

package service

import (
	"context"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/common/kafka"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/gateway/websocket"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
	"gorm.io/gorm"
)

type SendP2PMessageRequest struct {
	SenderID    uuid.UUID `json:"sender_id" binding:"required"`
	ReceiverID  uuid.UUID `json:"receiver_id" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	ContentType int       `json:"content_type"`
}

type SendGroupMessageRequest struct {
	SenderID    uuid.UUID `json:"sender_id" binding:"required"`
	GroupID     uuid.UUID `json:"group_id" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	ContentType int       `json:"content_type"`
}

type MessageService struct {
	DB            *gorm.DB
	KafkaProducer *kafka.Producer
}

func NewMessageService(db *gorm.DB, kafkaProducer *kafka.Producer) *MessageService {
	return &MessageService{
		DB:            db,
		KafkaProducer: kafkaProducer,
	}
}

func (m *MessageService) SendP2PMessage(ctx context.Context, req *SendP2PMessageRequest) (*websocket.MessageResponse, error) {
	// 1. 检查接收者是否是发送者的朋友
	var friendship types.Friends
	if err := m.DB.Where("user_id = ? AND friend_id = ?", req.SenderID, req.ReceiverID).First(&friendship).Error; err != nil {
		return nil, err
	}

	// 2. 创建消息struct
	message := types.P2PMessages{
		SenderID:    req.SenderID,
		Sender:      types.Users{ID: req.SenderID},
		ReceiverID:  req.ReceiverID,
		Receiver:    types.Users{ID: req.ReceiverID},
		Content:     req.Content,
		ContentType: req.ContentType,
		CreatedAt:   time.Now(),
	}
	// 3. 将消息发送到Kafka
	kafkaPayload := kafka.MessagePayload{
		Type:      "p2p_message",
		Data:      message,
		Timestamp: time.Now().Unix(),
	}
	if err := m.KafkaProducer.SendMessage(ctx, "p2p_message", "p2p_message", kafkaPayload); err != nil {
		slog.Error("Failed to send message to Kafka", "error", err)
	}
	// 4. 将消息存储到db.
	result := m.DB.Create(&message)
	if result.Error != nil {
		slog.Error("Failed to save message to database", "error", result.Error)
		return nil, result.Error
	}
	// 5. 返回消息结构
	return &websocket.MessageResponse{
		ID:        message.ID,
		Success:   true,
		Error:     "",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (m *MessageService) SendGroupMessage(ctx context.Context, req *SendGroupMessageRequest) (*websocket.MessageResponse, error) {
	// 1. 检查发送者是否是群成员
	var groupMember types.GroupMembers
	if err := m.DB.Where("user_id = ? AND group_id = ?", req.SenderID, req.GroupID).First(&groupMember).Error; err != nil {
		return nil, err
	}

	// 2. 创建消息struct
	groupMessage := types.GroupMessages{
		SenderID:    req.SenderID,
		Sender:      types.Users{ID: req.SenderID},
		GroupID:     req.GroupID,
		Group:       types.Groups{ID: req.GroupID},
		Content:     req.Content,
		ContentType: req.ContentType,
		CreatedAt:   time.Now(),
	}
	// 3. 将消息发送到Kafka
	kafkaPayload := kafka.MessagePayload{
		Type:      "group_message",
		Data:      groupMessage,
		Timestamp: time.Now().Unix(),
	}
	if err := m.KafkaProducer.SendMessage(ctx, "group_message", "group_message", kafkaPayload); err != nil {
		slog.Error("Failed to send message to Kafka", "error", err)
	}
	// 4. 将消息存储到db.
	result := m.DB.Create(&groupMessage)
	if result.Error != nil {
		slog.Error("Failed to save message to database", "error", result.Error)
		return nil, result.Error
	}
	// 5. 返回消息结构
	return &websocket.MessageResponse{
		ID:        groupMessage.ID,
		Success:   true,
		Error:     "",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (s *MessageService) GetP2PMessages(senderID, receiverID uuid.UUID, offset, limit int) ([]types.P2PMessages, error) {
	var messages []types.P2PMessages
	err := s.DB.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).
		Or("sender_id = ? AND receiver_id = ?", receiverID, senderID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (s *MessageService) GetGroupMessages(groupID uuid.UUID, offset, limit int) ([]types.GroupMessages, error) {
	var messages []types.GroupMessages
	err := s.DB.Where("group_id = ?", groupID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (s *MessageService) GetUserConversations(userID uuid.UUID) ([]types.Conversations, error) {
	var conversations []types.Conversations
	err := s.DB.Preload("LastMessage").Preload("Participants").
		Where("p2p_user1 = ? OR p2p_user2 = ?", userID, userID).
		Or("id IN (SELECT conversation_id FROM conversation_participants WHERE user_id = ?)", userID).
		Order("updated_at DESC").
		Find(&conversations).Error

	return conversations, err
}

// 更新P2P会话
func (s *MessageService) updateP2PConversation(message types.P2PMessages) {
	var conversation types.Conversations

	// 查找现有会话（确保user1 < user2的顺序）
	user1, user2 := message.SenderID, message.ReceiverID
	if user1.String() > user2.String() {
		user1, user2 = user2, user1
	}

	err := s.DB.Where("p2p_user1 = ? AND p2p_user2 = ?", user1, user2).First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新会话
		conversation = types.Conversations{
			P2PUser1:      user1,
			P2PUser2:      user2,
			LastMessageID: message.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		s.DB.Create(&conversation)

		// 创建参与者记录
		participants := []types.ConversationParticipants{
			{
				ConversationID: conversation.ID,
				UserID:         user1,
				UnreadCount:    0,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ConversationID: conversation.ID,
				UserID:         user2,
				UnreadCount:    1, // 接收者的未读数+1
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		s.DB.Create(&participants)
	} else {
		// 更新现有会话
		conversation.LastMessageID = message.ID
		conversation.UpdatedAt = time.Now()
		s.DB.Save(&conversation)

		// 更新接收者的未读计数
		s.DB.Model(&types.ConversationParticipants{}).
			Where("conversation_id = ? AND user_id = ?", conversation.ID, message.ReceiverID).
			Updates(map[string]interface{}{
				"unread_count": gorm.Expr("unread_count + 1"),
				"updated_at":   time.Now(),
			})
	}
}

func (s *MessageService) MarkMessagesAsRead(userID, conversationID uuid.UUID) error {
	return s.DB.Model(&types.ConversationParticipants{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{
			"unread_count": 0,
			"updated_at":   time.Now(),
		}).Error
}

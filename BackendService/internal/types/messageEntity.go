package types

import (
	"time"

	"github.com/google/uuid"
)

type P2PMessages struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	SenderID    uuid.UUID `gorm:"not null;column:sender_id;index"`
	Sender      Users
	ReceiverID  uuid.UUID `gorm:"not null;column:receiver_id;index"`
	Receiver    Users
	Content     string `gorm:"not null;column:content"`
	ContentType int    `gorm:"not null;column:content_type"`
	CreatedAt   time.Time
}

type GroupMessages struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	SenderID    uuid.UUID `gorm:"not null;column:sender_id;index"`
	Sender      Users
	Content     string    `gorm:"not null;column:content"`
	ContentType int       `gorm:"not null;column:content_type"`
	GroupID     uuid.UUID `gorm:"not null;column:group_id;index"`
	CreatedAt   time.Time
}

type Conversations struct {
	ID            uuid.UUID                  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	P2PUser1      uuid.UUID                  `gorm:"not null;column:p2p_user1_id;uniqueIndex:idx_p2p_user1_user2"`
	P2PUser2      uuid.UUID                  `gorm:"not null;column:p2p_user2_id;uniqueIndex:idx_p2p_user1_user2"`
	GroupID       uuid.UUID                  `gorm:"column:group_id;uniqueIndex:idx_group_id"`
	LastMessageID uuid.UUID                  `gorm:"column:last_message_id;index"`
	LastMessage   P2PMessages                `gorm:"foreignKey:LastMessageID;references:ID"`
	Participants  []ConversationParticipants `gorm:"many2many:group_members;joinForeignKey:ConversationID;joinReferences:UserID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ConversationParticipants struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	ConversationID uuid.UUID `gorm:"not null;column:conversation_id;uniqueIndex:idx_conversation_user"`
	UserID         uuid.UUID `gorm:"not null;column:user_id;uniqueIndex:idx_conversation_user"`
	UnreadCount    int       `gorm:"column:unread_count;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

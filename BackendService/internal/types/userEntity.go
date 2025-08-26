package types

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id" json:"id"`
	Username      string    `gorm:"unique;column:username" json:"username"`
	Nickname      string    `gorm:"column:nickname" json:"nickname"`
	ProfileAvatar string    `gorm:"column:profile_avatar" json:"profile_avatar"`
	Email         string    `gorm:"unique;column:email" json:"email"`
	Password      string    `gorm:"column:password" json:"-"`
	//OauthIdentities []OauthIdentity `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Friends                  []Users                    `gorm:"many2many:friends;joinForeignKey:UserID;joinReferences:FriendID"`
	Groups                   []Groups                   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE"`
	MemberGroups             []Groups                   `gorm:"many2many:group_members;joinForeignKey:UserID;joinReferences:GroupID"`
	ParticipantConversations []ConversationParticipants `gorm:"many2many:conversation_participants;joinForeignKey:UserID;joinReferences:ConversationID"`
	CreatedAt                time.Time                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                time.Time                  `gorm:"column:updated_at" json:"updated_at"`
}

type OauthIdentities struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	UserID     uuid.UUID `gorm:"not null;column:user_id;index"`
	User       Users     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Provider   string    `gorm:"not null;column:provider;uniqueIndex:idx_provider_identity"`
	ProviderID string    `gorm:"not null;column:provider_id;uniqueIndex:idx_provider_identity"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Friends struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	UserID    uuid.UUID `gorm:"not null;column:user_id;index"`
	FriendID  uuid.UUID `gorm:"not null;column:friend_id;index"`
	CreatedAt time.Time
}

type Groups struct {
	ID            uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	Name          string          `gorm:"not null;column:name"`
	OwnerID       uuid.UUID       `gorm:"not null;column:owner_id;index"`
	Members       []Users         `gorm:"many2many:group_members;joinForeignKey:GroupID;joinReferences:UserID"`
	GroupMessages []GroupMessages `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type GroupMembers struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	UserID    uuid.UUID `gorm:"not null;column:user_id;index"`
	GroupID   uuid.UUID `gorm:"not null;column:group_id;index"`
	Role      string    `gorm:"not null;column:role"`
	JoinedAt  time.Time
	UpdatedAt time.Time
}

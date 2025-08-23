package types

import "time"

type User struct {
	ID              string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id" json:"id"`
	Username        string          `gorm:"unique;column:username" json:"username"`
	Email           string          `gorm:"unique;column:email" json:"email"`
	Password        string          `gorm:"column:password" json:"-"`
	OauthIdentities []OauthIdentity `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Friendships     []User          `gorm:"many2many:friendships;"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

type OauthIdentity struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	UserID     string `gorm:"not null;column:user_id;index"`
	Provider   string `gorm:"not null;column:provider;uniqueIndex:idx_provider_identity"`
	ProviderID string `gorm:"not null;column:provider_id;uniqueIndex:idx_provider_identity"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Friendship struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();column:id"`
	UserID    string `gorm:"not null;column:user_id;index"`
	FriendID  string `gorm:"not null;column:friend_id;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

package user

import (
	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(user *types.User) error {
	return s.db.Create(user).Error
}

func (s *UserStore) GetUserByID(id string) (*types.User, error) {
	var user types.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

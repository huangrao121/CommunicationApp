package oauthauth

import "gorm.io/gorm"

type OauthStore struct {
	db *gorm.DB
}

func NewOauthStore(db *gorm.DB) *OauthStore {
	return &OauthStore{db: db}
}

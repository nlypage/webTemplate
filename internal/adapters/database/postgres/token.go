package postgres

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webTemplate/internal/domain/entity"
)

// tokenStorage is a struct that contains a pointer to a gorm.DB instance to interact with token repository.
type tokenStorage struct {
	db *gorm.DB
}

// NewTokenStorage is a function that returns a new instance of tokenStorage.
func NewTokenStorage(db *gorm.DB) *tokenStorage {
	return &tokenStorage{db: db}
}

// Create is a method to create a new Token in database.
func (s *tokenStorage) Create(ctx context.Context, token entity.Token) (*entity.Token, error) {
	err := s.db.WithContext(ctx).Create(&token).Error
	return &token, err
}

// GetByUserID is a method that returns an error and a pointer to a Token instance by user id and token type.
func (s *tokenStorage) GetByUserID(ctx context.Context, userID string, tokenType string) (*entity.Token, error) {
	var token *entity.Token
	err := s.db.WithContext(ctx).Model(&entity.Token{}).Where(
		"id = ? AND type = ? AND expires > ?", userID, tokenType, time.Now(),
	).First(&token).Error
	return token, err
}

// DeleteAll is a method to delete all user Tokens in database.
func (s *tokenStorage) DeleteAll(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Delete(&entity.Token{}, "user_id = ?", userID).Error
}

// Delete is a method to delete an existing Token in database by user id and token type.
func (s *tokenStorage) Delete(ctx context.Context, userID string, tokenType string) error {
	return s.db.WithContext(ctx).Delete(&entity.Token{}, "user_id = ? AND type = ?", userID, tokenType).Error
}

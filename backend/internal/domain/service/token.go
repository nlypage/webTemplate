package service

import (
	"context"
	"github.com/spf13/viper"
	"time"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/entity"
	"webTemplate/internal/domain/utils/auth"
)

type TokenStorage interface {
	Create(ctx context.Context, token entity.Token) (*entity.Token, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Token, error)
	DeleteAll(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string, tokenType string) error
}

// tokenService is a struct that contains a pointer to a gorm.DB instance to interact with token repository.
type tokenService struct {
	storage TokenStorage
}

func NewTokenService(storage TokenStorage) *tokenService {
	return &tokenService{storage: storage}
}

// GenerateToken is a method to generate a new token.
func (s *tokenService) GenerateToken(ctx context.Context, userID string, expires time.Time, tokenType string) (*entity.Token, error) {
	jwtToken, err := auth.GenerateToken(userID, expires, tokenType)
	if err != nil {
		return nil, err
	}

	token, err := s.storage.Create(ctx, entity.Token{
		Token:   jwtToken,
		UserID:  userID,
		Type:    tokenType,
		Expires: expires,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// DeleteToken is a method to delete a token by user id and token type.
func (s *tokenService) DeleteToken(ctx context.Context, userID string, tokenType string) error {
	return s.storage.Delete(ctx, userID, tokenType)
}

// GenerateAuthTokens is a method to generate access and refresh tokens.
func (s *tokenService) GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error) {
	authToken, err := s.GenerateToken(
		c,
		userID,
		time.Now().UTC().Add(time.Minute*time.Duration(viper.GetInt("service.backend.jwt.access-token-expiration"))),
		auth.TokenTypeAccess,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateToken(
		c,
		userID,
		time.Now().UTC().Add(time.Minute*time.Duration(viper.GetInt("service.backend.jwt.refresh-token-expiration"))),
		auth.TokenTypeRefresh,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokens{
		Access: dto.Token{
			Token:   authToken.Token,
			Expires: authToken.Expires,
		},
		Refresh: dto.Token{
			Token:   refreshToken.Token,
			Expires: refreshToken.Expires,
		},
	}, nil
}

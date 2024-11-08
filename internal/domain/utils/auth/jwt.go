package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"strings"
	"time"
	"webTemplate/internal/domain/common/errorz"
	"webTemplate/internal/domain/entity"
)

func VerifyToken(authHeader, secret, tokenType string) (string, error) {
	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenStr == "" {
		return "", errorz.AuthHeaderIsEmpty
	}

	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	jwtType, ok := claims["type"].(string)
	if !ok || jwtType != tokenType {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token sub")
	}

	return userID, nil
}

func GetUserFromJWT(jwt, tokenType string, context context.Context, getUser func(context.Context, string) (*entity.User, error)) (*entity.User, error) {
	id, errVerify := VerifyToken(jwt, viper.GetString("service.backend.jwt.secret"), tokenType)
	if errVerify != nil {
		return &entity.User{}, errVerify
	}

	user, errGetUser := getUser(context, id)
	if errGetUser != nil {
		return &entity.User{}, errGetUser
	}

	return user, nil
}

func GenerateToken(userID string, expires time.Time, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(viper.GetString("service.backend.jwt.secret")))
}

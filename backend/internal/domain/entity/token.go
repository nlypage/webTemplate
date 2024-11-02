package entity

import "time"

// Token is a struct that represents an authorization token in database.
type Token struct {
	ID        string `gorm:"primaryKey;not null;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Token   string    `gorm:"not null"`
	UserID  string    `gorm:"not null;type:uuid"`
	Type    string    `gorm:"not null"`
	Expires time.Time `gorm:"not null"`
	User    *User     `gorm:"foreignKey:user_id;references:id"`
}

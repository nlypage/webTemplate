package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is a struct that represents a user in database.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;not null;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Email         string  `json:"email" gorm:"uniqueIndex"`
	VerifiedEmail bool    `json:"verified_email" gorm:"default:false;not null"`
	Password      []byte  `json:"-"`
	Role          string  `json:"role" gorm:"default:user;not null"`
	Token         []Token `json:"-" gorm:"foreignKey:user_id;references:id"`
	Username      string  `json:"username"`
}

// HashedPassword is a function to hash the password.
func HashedPassword(password string) []byte {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return hashedPassword
}

// SetPassword is a method to hash the password before storing it.
func (user *User) SetPassword(password string) {
	user.Password = HashedPassword(password)
}

// ComparePassword is a method to compare the password with the hashed password.
func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

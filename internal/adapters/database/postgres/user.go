package postgres

import (
	"context"
	"gorm.io/gorm"
	"webTemplate/internal/domain/entity"
)

// userStorage is a struct that contains a pointer to a gorm.DB instance to interact with user repository.
type userStorage struct {
	db *gorm.DB
}

// NewUserStorage is a function that returns a new instance of usersStorage.
func NewUserStorage(db *gorm.DB) *userStorage {
	return &userStorage{db: db}
}

// Create is a method to create a new User in database.
func (s *userStorage) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	err := s.db.WithContext(ctx).Create(&user).Error
	return &user, err
}

// GetByID is a method that returns an error and a pointer to a User instance by id.
func (s *userStorage) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user *entity.User
	err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).First(&user).Error
	return user, err
}

// GetAll is a method that returns a slice of pointers to all User instances.
func (s *userStorage) GetAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	var users []entity.User
	err := s.db.WithContext(ctx).Model(&entity.User{}).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Update is a method to update an existing User in database.
func (s *userStorage) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(&user).Error
	return user, err
}

// Delete is a method to delete an existing User in database.
func (s *userStorage) Delete(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&entity.User{}, "id = ?", id).Error
}

// GetByEmail is a method that returns a pointer to a User instance and error by email.
func (s *userStorage) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

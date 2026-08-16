package auth

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	UpdateProfile(userID uint, username, email string) error
	UpdatePassword(userID uint, hashedPassword string) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	CreateToken(token *Token) error
	DeleteTokensByUserID(userID uint, tokenType string) error
	DeleteToken(tokenStr string, tokenType string) error
	FindByToken(tokenStr string) (*Token, error)
	DeleteExpiredTokens() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) UpdateProfile(userID uint, username, email string) error {
	return r.db.Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"username": username,
			"email":    email,
		}).Error
}

func (r *repository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Preload("Role").Preload("Center").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Role").Preload("Center").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateToken(token *Token) error {
	return r.db.Create(token).Error
}

func (r *repository) DeleteTokensByUserID(userID uint, tokenType string) error {
	return r.db.Where("user_id = ? AND type = ?", userID, tokenType).
		Delete(&Token{}).Error
}

func (r *repository) DeleteToken(tokenStr string, tokenType string) error {
	return r.db.Where("token = ? AND type = ?", tokenStr, tokenType).
		Delete(&Token{}).Error
}

func (r *repository) FindByToken(tokenStr string) (*Token, error) {
	var token Token
	err := r.db.Where("token = ? AND expires_at > ?", tokenStr, time.Now()).
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *repository) DeleteExpiredTokens() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&Token{})
	log.Printf("Maintenance: Deleted %d expired tokens", result.RowsAffected)
	return result.Error
}

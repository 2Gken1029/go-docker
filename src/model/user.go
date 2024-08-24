package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required,min=1,max=100" ja:"名前"`
	Email     string    `json:"email"  validate:"uniqueEmail,required,email" ja:"メールアドレス"`
	Password  string    `json:"password" validate:"required,min=8,max=100" ja:"パスワード"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUserRequest struct {
	Name     string      `json:"name" validate:"required,min=1,max=100" ja:"名前"`
	Email    string      `json:"email"  validate:"uniqueEmail,required,email" ja:"メールアドレス"`
	Password string      `json:"password" validate:"required,min=8,max=100" ja:"パスワード"`
}

type LoginRequest struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

/*
新しいユーザーを作成。

@args db *gorm.DB
@args request *RegisterUserRequest
*/
func CreateUser(db *gorm.DB, request *RegisterUserRequest) (uint, error) {
	// パスワードをハッシュ化
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return 0, err
	}

	// ユーザー情報を構造体にセット
	user := &User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	// データベースにユーザーを作成
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}

	// 作成されたユーザの主キー(ID)を取得
	userID := user.ID

	return userID, nil
}

/*
パスワードをハッシュ化。

@args password string
@return string(ハッシュ化したパスワード), error
*/
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

/*
指定されたメールアドレスからユーザ情報を取得。
@args db *gorm.DB
@args email string
@return User（ユーザ情報）, error
*/
func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := User{}
	err := db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	return &user, nil
}
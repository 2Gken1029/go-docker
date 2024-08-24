package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type InvalidToken struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"not null" ja:"リフレッシュトークン"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null" ja:"期限"`
}

// 無効にするリフレッシュトークンを登録します
func CreateInvalidToken(db *gorm.DB, token string, expireAte float64) error {

	intExpireAte := int64(expireAte)
	t := time.Unix(intExpireAte, 0)

	var invalidToken *InvalidToken
	createdAte := time.Now()
	invalidToken = &InvalidToken{
		Token:     token,
		CreatedAt: createdAte,
		ExpiredAt: t,
	}

	if err := db.Create(&invalidToken).Error; err != nil {
		return err
	}

	return nil
}

// 無効なトークン情報を取得
func GetInvalidTokenByToken(db *gorm.DB, token string) (*InvalidToken, error) {
	invalidToken := InvalidToken{}
	err := db.Where("token = ?", token).First(&invalidToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	return &invalidToken, nil
}

// 期限切れ無効化トークンの物理削除
func DeleteExpiredTokens(db *gorm.DB) (int64, error) {
	now := time.Now().UTC()
	result := db.Unscoped().Where("expired_at < ?", now).Delete(&InvalidToken{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

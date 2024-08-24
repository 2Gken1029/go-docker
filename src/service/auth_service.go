package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

/*
JWTの作成

@args email string
@args password string
@args user_id string
@return access_token, refresh_token
*/
func (auth *AuthService) CreateJWT(email string, name string, user_id int) (TokenResponse, error) {
	// JWTに付与する構造体
	claims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"iss":   "eco_app",
		"sub":   user_id,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 72時間が有効期限
	}

	// ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名を付与
	accessToken, _ := token.SignedString([]byte(os.Getenv("ACCESS_SECRET_KEY")))
	fmt.Println("accessToken:", accessToken)

	// JWTに付与する構造体（リフレッシュトークン）
	refreshClaims := jwt.MapClaims{
		"email": email,
		"name":  name,
		"iss":   "eco_app",
		"sub":   user_id,
		"exp":   time.Now().Add(time.Hour * 24 * 90).Unix(), // 90日が有効期限
	}

	// ヘッダーとペイロード生成（リフレッシュトークン）
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, _ := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))

	// TokenResponse構造体にアクセストークンとリフレッシュトークンを格納
	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}

	return response, nil
}

package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// ParseAndValidateAccessToken は与えられたコンテキストからアクセストークンを解析し、検証する。
//
// Parameters:
//   - c: echo.Context - Echoのコンテキスト
//
// Returns:
//   - *jwt.Token: トークンが有効な場合、解析されたJWTトークン
//   - error: エラーが発生した場合、そのエラー
//
// この関数は Authorization ヘッダーからトークンを取得し、
// それを検証する。検証に成功した場合、解析されたトークンを返す。
// トークンが無効な場合やエラーが発生した場合、適切なエラーレスポンスを生成して返す。
func ParseAndValidateAccessToken(c echo.Context) (*jwt.Token, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, SendErrorResponse(c, http.StatusUnauthorized, "Missing Access Token")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET_KEY")), nil
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if token == nil || !token.Valid {
		return nil, SendErrorResponse(c, http.StatusInternalServerError, "認証できません")
	}

	return token, nil
}

// ParseAndValidateSecretToken は与えられたコンテキストからシークレットトークンを解析し、検証する。
//
// Parameters:
//   - c: echo.Context - Echoのコンテキスト
//
// Returns:
//   - *jwt.Token: トークンが有効な場合、解析されたJWTトークン
//   - string: トークンが有効な場合、現在のリフレッシュトークン
//   - error: エラーが発生した場合、そのエラー
//
// この関数は Authorization ヘッダーからトークンを取得し、
// それを検証する。検証に成功した場合、解析されたトークンを返す。
// トークンが無効な場合やエラーが発生した場合、適切なエラーレスポンスを生成して返す。
func ParseAndValidateSecretToken(c echo.Context) (*jwt.Token, string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, "", SendErrorResponse(c, http.StatusUnauthorized, "Missing Refresh Token")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SECRET_KEY")), nil
	})
	if err != nil {
		log.Printf(err.Error())
		return nil, "", SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if token == nil || !token.Valid {
		return nil, "", SendErrorResponse(c, http.StatusInternalServerError, "Forbidden: Invalid Claims")
	}

	return token, tokenString, nil
}

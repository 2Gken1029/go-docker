package controller

import (
	"echo-get-started/model"
	"echo-get-started/service"
	"echo-get-started/util"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserController interface{
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type userController struct {
	Database *gorm.DB
}

func NewUserController(db *gorm.DB) IUserController {
	return &userController{
		Database: db,
	}
}

/*
新規登録

@args c: echo.Context (名前とメールアドレスとパスワード)
*/
func (uc *userController) Register(c echo.Context) error {
	request := &model.RegisterUserRequest{}

	// バリデーション
	if err := c.Bind(request); err != nil {
		return util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(request); err != nil {
		return util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// DBインサート
	userId, err := model.CreateUser(uc.Database, request)
	if err != nil && userId == 0 {
		log.Printf("データベース作成（user）に失敗しました : " + err.Error())
		return fmt.Errorf(err.Error())
	}

	// JWT発行
	auth_service := service.AuthService{}
	token, err := auth_service.CreateJWT(request.Email, request.Name, int(userId))
	if err != nil {
		log.Printf("JWT作成に失敗しました : " + err.Error())
		return util.SendErrorResponse(c, http.StatusInternalServerError, "JWT作成に失敗しました")
	}
	return util.SendSuccessResponse(c, token)
}

/*
ログイン

@args c: echo.Context (名前とメールアドレスとパスワード)
*/
func (uc *userController) Login(c echo.Context) error {
	request := &model.LoginRequest{}

	if err := c.Bind(request); err != nil {
		return util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(request); err != nil {
		return util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// DBからユーザー情報を取得する
	user, err := model.GetUserByEmail(uc.Database, request.Email)
	if err != nil && user == nil {
		return util.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect Email Address or Password")
	} else if err != nil {
		return util.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}

	// データベースから取得したハッシュとパスワードを比較
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return util.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect Email Address or Password")
	} else if err != nil {
		return util.SendErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}

	// OKなら、JWTを発行する
	auth_service := service.AuthService{}
	token, _ := auth_service.CreateJWT(user.Email, user.Name, int(user.ID))
	return util.SendSuccessResponse(c, token)
}


/*
リフレッシュ
リフレッシュトークンを使用して、アクセストークンとリフレッシュトークンを再発行する

@args c: echo.Context (名前とメールアドレスとパスワード)
*/
func (uc *userController) Refresh(c echo.Context) error {
	// JWT検証
	token, tokenString, err := util.ParseAndValidateSecretToken(c)
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// DBからユーザー情報を取得する
		user, err := model.GetUserByEmail(uc.Database, claims["email"].(string))
		if err != nil && user == nil {
			return util.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		}

		// InvalidTokenが無効化されているかチェックする
		_, err = model.GetInvalidTokenByToken(uc.Database, string(tokenString))
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return util.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		} else {
			// 無効化されている場合
			return util.SendErrorResponse(c, http.StatusUnauthorized, "Invalid Token")
		}

		// JWTを作成して、新トークンを返す
		auth_service := service.AuthService{}
		token, _ := auth_service.CreateJWT(user.Email, user.Name, int(user.ID))
		return util.SendSuccessResponse(c, token)

	} else {
		return util.SendErrorResponse(c, http.StatusBadRequest, "認証できません")
	}
}
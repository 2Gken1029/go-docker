package util

import (
	"echo-get-started/model"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CustomValidator struct {
	trans     ut.Translator
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	if err == nil {
		return err
	}

	errs := err.(validator.ValidationErrors)
	msg := ""
	for _, ve := range errs.Translate(cv.trans) {
		if msg != "" {
			msg += ", "
		}
		msg += ve
	}
	return errors.New(msg)
}

func NewValidator() echo.Validator {
	english := en_US.New() // English translator
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	//メールアドレスの重複チェックのカスタムバリデーション
	validate.RegisterValidation("uniqueEmail", ValidateUniqueEmail)
	validate.RegisterTranslation("uniqueEmail", trans,
		func(ut ut.Translator) error {
			return ut.Add("uniqueEmail", "{0}はご利用できないです。", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("uniqueEmail", fe.Field())
			return t
		})

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Fatal(err)
	}
	return &CustomValidator{
		trans:     trans,
		validator: validate,
	}
}

// バリデーション用関数
func ValidateUniqueEmail(fl validator.FieldLevel) bool {
	//dbに問い合わせて既に存在していればバリデーションエラーにする
	connectInfo := fmt.Sprintf(
		"%s:%s@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(connectInfo), &gorm.Config{}) // dbとerrを宣言および初期化
	if err != nil {
		panic(err.Error())
	}

	var user []model.User
	result := db.Where("email = ?", fl.Field().String()).Limit(1).Find(&user)
	return result.RowsAffected == 0
}
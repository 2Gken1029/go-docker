package config

import (
	"echo-get-started/controller"
	"echo-get-started/util"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ルーティング
type Routing struct {
	DB   *DB
	Port string
}

func NewRouting(db *DB) *Routing {
	c := NewConfig()
	r := &Routing{
		DB:   db,
		Port: c.Routing.Port,
	}
	r.startRouting()
	return r
}

func (r *Routing) startRouting() {
	e := echo.New()
	e.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	e.Logger.SetOutput(os.Stderr)

	// バリデーションの設定
	e.Validator = util.NewValidator()

	// 共通ミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// コントローラー定義
	userController := controller.NewUserController(r.DB.Connect())

	// ELBのヘルスチェック用
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "health check") })

	v1 := e.Group("/v1")
	{
		v1.POST("/login", userController.Login)

		userGroup := v1.Group("/users")
		{
			userGroup.POST("", userController.Register)
		}
	}

	e.Logger.Fatal(e.Start(r.Port))
}

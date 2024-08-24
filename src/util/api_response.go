package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/*
成功時のレスポンス

@args c: echo.Context
@args data: 返却データ
*/
func SendSuccessResponse(c echo.Context, data interface{}, message ...string) error {
	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}

	// message パラメータが渡された場合は response に含める
	if len(message) > 0 {
		response["message"] = message[0]
	}

	return c.JSON(http.StatusOK, response)

}

/*
エラー時のレスポンス

@args c: echo.Context
@args code: レスポンスコード
@args message: エラーメッセージ
*/
func SendErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

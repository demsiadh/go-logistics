package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response gin的统一返回格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

func SuccessResponseWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func (r *Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type ErrorMsg struct {
	Code    int
	Message string
}

var (
	UserNotFound  = &ErrorMsg{Code: 10001, Message: "用户不存在"}
	UserInfoError = &ErrorMsg{Code: 10002, Message: "用户信息错误"}
	ServerError   = &ErrorMsg{Code: 50001, Message: "服务器错误"}
)

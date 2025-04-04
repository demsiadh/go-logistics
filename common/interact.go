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

func ErrorResponse(c *gin.Context, err *ErrorMsg) {
	c.JSON(http.StatusOK, Response{
		Code:    err.Code,
		Message: err.Message,
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

func (e *ErrorMsg) Error() string {
	return e.Message
}

var (
	ServerError    = &ErrorMsg{Code: 50001, Message: "服务器错误"}
	ParamError     = &ErrorMsg{Code: 40001, Message: "参数错误"}
	RecordNotFound = &ErrorMsg{Code: 60001, Message: "记录不存在"}
	RecordExist    = &ErrorMsg{Code: 60002, Message: "记录已存在"}
)

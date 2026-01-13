package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应结构体
type Response struct {
	// Code 业务状态码
	// 200 表示成功，其他值表示具体的业务错误
	Code int `json:"code"`

	// Msg 友好提示信息
	// 用于前端展示给用户的消息
	Msg string `json:"msg"`

	// ErrMsg 错误详情
	// 仅在调试模式或特定异常时返回，用于排查问题
	// omitempty 表示如果为空则不序列化该字段
	ErrMsg string `json:"errMsg,omitempty"`

	// Data 业务数据
	// 存放实际返回的业务数据对象或列表
	Data interface{} `json:"data"`

	// Timestamp 响应生成时间戳
	// 格式：ISO8601 (2006-01-02T15:04:05Z07:00)
	Timestamp string `json:"timestamp"`
}

// newResponse 创建一个新的响应对象
func newResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code:      code,
		Msg:       msg,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// Success 成功响应
//
// 参数:
//
//	c: Gin 上下文
//	data: 返回的业务数据
func Success[T any](c *gin.Context, data T) {
	resp := newResponse(200, "操作成功", data)
	c.JSON(http.StatusOK, resp)
}

// Error 错误响应
//
// 参数:
//
//	c: Gin 上下文
//	code: 错误码
//	msg: 错误提示信息 (如果为空，尝试根据 code 获取国际化消息)
func Error(c *gin.Context, code int, msg string) {
	if msg == "" {
		msg = GetMsg(code, c.GetHeader("Accept-Language"))
	}
	resp := newResponse(code, msg, nil)
	c.JSON(http.StatusOK, resp) // 业务错误通常也返回 200 HTTP 状态码，通过 Body 中的 Code 区分
}

// ErrorWithDetail 带详情的错误响应
//
// 参数:
//
//	c: Gin 上下文
//	code: 错误码
//	msg: 错误提示信息
//	errMsg: 详细错误信息 (仅供开发者查看)
func ErrorWithDetail(c *gin.Context, code int, msg string, errMsg string) {
	if msg == "" {
		msg = GetMsg(code, c.GetHeader("Accept-Language"))
	}
	resp := newResponse(code, msg, nil)
	resp.ErrMsg = errMsg
	c.JSON(http.StatusOK, resp)
}

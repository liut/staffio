package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ok ...
func Ok(c *gin.Context, args ...interface{}) {
	Out(c, 200, args...)
}

// Out ...
func Out(c *gin.Context, code int, args ...interface{}) {
	res := gin.H{"status": 0, "ok": true} // status is deprecated
	if len(args) > 0 {
		res["data"] = args[0]
		if len(args) > 1 {
			res["total"] = args[1]
		}
	}
	c.JSON(code, res)
}

// Fail response fail with code, err, field
func Fail(c *gin.Context, code int, args ...interface{}) {
	if len(args) == 0 {
		c.AbortWithStatus(code)
		return
	}
	var res RespFail
	res.Error = GetError(c.Request, code, args[0], args[1:]...)
	c.AbortWithStatusJSON(code, res)
}

// RespDone 操作成功返回的结构
type RespDone struct {
	Ok    bool        `json:"ok" description:"操作成功"` // OK
	Data  interface{} `json:"data,omitempty"`        // main data
	Extra interface{} `json:"extra,omitempty"`       // extra data
}

// RespFail 出现错误，返回相关的错误码和消息文本
type RespFail struct {
	Ok    bool  `json:"ok" description:"操作失败"`
	Error Error `json:"error" description:"错误集"`
}

// IError ...
type IError interface {
	GetCode() int
	GetMessage() string
	GetField() string
}

// Error ...
type Error struct {
	Code    int    `json:"code" description:"错误代码"`
	Message string `json:"message" description:"错误信息"`
	Field   string `json:"field,omitempty" description:"错误字段,可选,多用于表单校验"`
}

// ICodeErrorReq ...
type ICodeErrorReq interface {
	Code() int
	ErrorReq(r *http.Request) string
}

// IFieldErrorReq ...
type IFieldErrorReq interface {
	ErrorReq(r *http.Request) string
	Field() string
}

func GetError(r *http.Request, code int, err interface{}, args ...interface{}) Error {
	var field string
	if len(args) > 0 {
		if v, ok := args[0].(string); ok {
			field = v
		}
	}
	switch e := err.(type) {
	case Error:
		e.Field = field
		return e
	case *Error:
		e.Field = field
		return *e
	case ICodeErrorReq:
		return Error{Code: e.Code(), Message: e.ErrorReq(r), Field: field}
	case IFieldErrorReq:
		return Error{Code: code, Message: e.ErrorReq(r), Field: e.Field()}
	case string:
		return Error{Code: code, Message: e, Field: field}
	case error:
		return Error{Code: code, Message: e.Error(), Field: field}
	case interface{ GetMessage() string }:
		return Error{Code: code, Message: e.GetMessage(), Field: field}
	default:
		if code >= 100 && code < 600 {
			return Error{Code: code, Message: http.StatusText(code), Field: field}
		}
		return Error{Code: code, Message: "unkown error", Field: field}
	}
}

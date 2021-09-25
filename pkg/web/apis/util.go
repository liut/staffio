package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/web/i18n"
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
	res := gin.H{"status": code, "ok": false} // status is deprecated
	if len(args) > 0 {
		var msg string
		if str, ok := args[0].(string); ok {
			msg = str
		} else if err, ok := args[0].(error); ok {
			msg = err.Error()
		} else {
			msg = http.StatusText(code)
		}
		res["message"] = msg
		if len(args) > 1 {
			res["field"] = args[1]
		}
	}

	c.AbortWithStatusJSON(code, res)
}

type respOK struct {
	Ok bool `json:"ok,required" description:"操作成功"`
}

// respDone 操作成功返回的结构
type respDone struct {
	Ok    bool        `json:"ok,required" description:"操作成功"`
	Data  interface{} `json:"data,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

// 出现错误，返回相关的错误码和消息文本
type respError struct {
	Ok     bool    `json:"ok" description:"操作失败"`
	Code   int     `json:"code" description:"错误代码"`
	Errors []Error `json:"errors" description:"错误集"`
}

func (re *respError) GetCode() int {
	return re.Code
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
	Message string `json:"message,required" description:"错误信息"`
	Field   string `json:"field,omitempty" description:"错误字段,可选,多用于表单校验"`
}

// FieldError ...
type FieldError interface {
	Field() string
}

func getErrors(r *http.Request, code int, err interface{}, args ...string) (errors []Error) {
	var field string
	if len(args) > 0 {
		field = args[0]
	}
	switch e := err.(type) {
	case Error:
		return append(errors, Error{Code: e.Code, Message: e.Message, Field: e.Field})
	case *Error:
		return append(errors, Error{Code: e.Code, Message: e.Message, Field: e.Field})
	case i18n.ErrorValue:
		return append(errors, Error{Code: e.Code(), Message: e.ErrorString(i18n.GetPrinter(r)), Field: field})
	case FieldError:
		return append(errors, Error{Message: i18n.GetFieldErrorString(i18n.GetPrinter(r), e), Field: e.Field()})
	case []FieldError:
		for _, _e := range e {
			errors = append(errors, Error{Message: i18n.GetFieldErrorString(i18n.GetPrinter(r), _e), Field: _e.Field()})
		}
		return
	case string:
		return append(errors, Error{Code: code, Message: e, Field: field})
	case interface{ GetMessage() string }:
		return append(errors, Error{Code: code, Message: e.GetMessage(), Field: field})
	case error:
		return append(errors, Error{Code: code, Message: e.Error(), Field: field})
	}

	return
}

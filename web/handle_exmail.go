package web

import (
	"github.com/RangelReale/osin"
	"github.com/wealthworks/go-tencent-api/exmail"

	. "lcgc/platform/staffio/settings"
)

func countNewMail(ctx *Context) (err error) {
	if !ctx.checkLogin() {
		return nil
	}
	email := ctx.User.Uid + "@" + Settings.EmailDomain
	count, err := exmail.CountNewMail(email)
	res := make(osin.ResponseData)
	res["email"] = email
	res["unseen"] = count
	return outputJson(res, ctx.Writer)
}

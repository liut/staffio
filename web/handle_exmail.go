package web

import (
	"encoding/binary"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/wealthworks/go-tencent-api/exmail"
	"log"
	"net/http"

	. "lcgc/platform/staffio/settings"
)

func countNewMail(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	email := ctx.User.Uid + "@" + Settings.EmailDomain
	res := make(osin.ResponseData)
	res["email"] = email
	key := []byte(fmt.Sprintf("mail-count-%s", ctx.User.Uid))

	if bv, err := cache.Get(key); err == nil {
		res["unseen"] = binary.LittleEndian.Uint32(bv)
	} else {
		count, err := exmail.CountNewMail(email)
		if err != nil {
			log.Printf("check new mail failed: %s", err)
			return nil
		}
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(count))
		cache.Set(key, bs, int(Settings.CacheLifetime))
		res["unseen"] = count
		res["got"] = true
	}

	return outputJson(res, ctx.Writer)
}

func loginToExmail(ctx *Context) error {
	if !ctx.checkLogin() {
		return nil
	}
	email := ctx.User.Uid + "@" + Settings.EmailDomain
	url, err := exmail.GetLoginURL(email)
	if err != nil {
		return err
	}
	http.Redirect(ctx.Writer, ctx.Request, url, http.StatusFound)
	return nil
}

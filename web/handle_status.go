package web

import (
	"fmt"
	"net/http"
	"tuluu.com/liut/keeper"
	. "tuluu.com/liut/staffio/settings"
)

func handleStatus(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil || !ctx.User.IsKeeper() {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}

	keeper.BootstrapPrefix = fmt.Sprintf("%sbootstrap-3.3.5/", Settings.ResUrl)

	switch ctx.Vars["topic"] {
	case "monitor":
		keeper.HandleMonitor(w, req)
	case "stacks":
		keeper.HandleStack(w, req)
	default:
		http.NotFound(w, req)
	}
	return nil
}

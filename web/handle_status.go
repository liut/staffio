package web

import (
	"fmt"
	"net/http"

	"lcgc/liut/keeper"
	. "lcgc/platform/staffio/settings"
)

func handleStatus(w http.ResponseWriter, req *http.Request, ctx *Context) (err error) {
	if ctx.User == nil || !ctx.User.IsKeeper() {
		http.Redirect(w, req, reverse("login"), http.StatusTemporaryRedirect)
		return nil
	}

	keeper.BootstrapPrefix = fmt.Sprintf("%sbootstrap-3.3.5/", Settings.ResUrl)

	switch ctx.Vars["topic"] {
	case "monitor":
		return T("dust_status.html").Execute(w, map[string]interface{}{
			"SysStatus": keeper.CurrentSystemStatus(),
			"ctx":       ctx,
		})
	case "stacks":
		keeper.HandleStack(w, req)
	default:
		http.NotFound(w, req)
	}
	return nil
}

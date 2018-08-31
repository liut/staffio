package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liut/keeper"
)

func (s *server) handleStatus(c *gin.Context) {
	req := c.Request
	w := c.Writer

	switch c.Param("topic") {
	case "monitor":
		s.Render(c, "dust_status.html", map[string]interface{}{
			"SysStatus": keeper.CurrentSystemStatus(),
			"ctx":       c,
		})
	case "stacks":
		keeper.HandleStack(w, req)
	default:
		http.NotFound(w, req)
	}
}

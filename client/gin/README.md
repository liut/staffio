Staffio client for gin
===


Example
---

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	staffiogin "github.com/liut/staffio/client/gin"
)

func main() {

	router := gin.Default()
	loginPath := "/auth/login"
	staffiogin.SetLoginPath(loginPath)
	staffiogin.SetAdminPath("/admin")

	router.GET(loginPath, staffiogin.LoginHandler)
	router.GET("/auth/callback", staffiogin.AuthCodeCallback("admin"))

	adminGroup := router.Group("/admin", staffiogin.AuthMiddleware(true)) // auto redirect
	adminGroup.GET("/", handlerAdminWelcome)
	...

	apiGroup := router.Group("/api", AuthMiddleware(false)) // don't redirect
	apiGroup.GET("/me", staffiogin.HandlerShowMe)
	...

}

func handlerAdminWelcome(c *gin.Context) {
	user, _ := staffiogin.UserWithContext(c)
	...
}

```

Staffio client
===


Example
---

```go
package main

import (
	"fmt"
	"net/http"

	staffio "github.com/liut/staffio/client"
)

func main() {

	loginPath := "/auth/login"
	staffio.LoginPath = loginPath
	staffio.AdminPath = "/admin"

	http.HandleFunc(loginPath, staffio.LoginHandler)
	http.Handle("/auth/callback", staffio.AuthCodeCallback("admin"))

	authF1 := staffio.AuthMiddleware(true) // auto redirect
	http.Handle("/admin", authF1(http.HandlerFunc(handlerAdminWelcome)))
	// more handlers
}

func handlerAdminWelcome(w http.ResponseWriter, r *http.Request) {
	user := staffio.UserFromContext(r.Context())
	fmt.Printf("user: %s", user.Name)
}

```

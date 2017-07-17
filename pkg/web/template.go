package web

import (
	"html/template"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"

	"lcgc/platform/staffio/pkg/settings"
)

var (
	cachedTemplates = map[string]*template.Template{}
	cachedMutex     sync.Mutex
	// funcs           = template.FuncMap{
	// 	"reverse": reverse,
	// }
)

// func reverse(name string, things ...interface{}) string {
// 	//convert the things to strings
// 	strs := make([]string, len(things))
// 	for i, th := range things {
// 		strs[i] = fmt.Sprint(th)
// 	}
// 	//grab the route
// 	u, err := ws.GetRoute(name).URL(strs...)
// 	if err != nil {
// 		log.Printf("GetRoute err %s", err)
// 		return "/" + name
// 		// panic(err)
// 	}
// 	return u.Path
// }

const (
	kReferer = "_ref"
)

func refererWithContext(c *gin.Context) (ref string) {
	cookie, err := c.Cookie(kReferer)
	if err == nil {
		ref = cookie
	}
	return
}

func markReferer(c *gin.Context) {
	c.SetCookie(kReferer, c.Request.RequestURI, 10, "/", "", false, true)
}

func Render(c *gin.Context, name string, data interface{}) (err error) {
	instance := T(name)
	if m, ok := data.(map[string]interface{}); ok {
		m["appVersion"] = settings.Version()
		m["navSimple"] = false
		session := ginSession(c)
		m["session"] = session
		m["referer"] = refererWithContext(c)
		var user *User
		v, exist := c.Get(kAuthUser)
		if exist {
			user = v.(*User)
		} else {
			user, err = UserFromRequest(c.Request)
		}
		m["currUser"] = user
		err = instance.Execute(c.Writer, m)
	} else {
		err = instance.Execute(c.Writer, data)
	}
	return
}

func T(name string) *template.Template {
	cachedMutex.Lock()
	defer cachedMutex.Unlock()

	if t, ok := cachedTemplates[name]; ok {
		return t
	}

	t := template.New("_base.html")
	t = template.Must(t.ParseFiles(
		filepath.Join(settings.Root, "templates/_base.html"),
		filepath.Join(settings.Root, "templates", name),
	))
	cachedTemplates[name] = t

	return t
}

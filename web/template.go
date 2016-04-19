package web

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sync"

	. "lcgc/platform/staffio/settings"
)

var (
	cachedTemplates = map[string]*template.Template{}
	cachedMutex     sync.Mutex
	funcs           = template.FuncMap{
		"reverse": reverse,
	}
)

func reverse(name string, things ...interface{}) string {
	//convert the things to strings
	strs := make([]string, len(things))
	for i, th := range things {
		strs[i] = fmt.Sprint(th)
	}
	//grab the route
	u, err := router.GetRoute(name).URL(strs...)
	if err != nil {
		return "/" + name
		// panic(err)
	}
	return u.Path
}

func T(name string) *template.Template {
	cachedMutex.Lock()
	defer cachedMutex.Unlock()

	if t, ok := cachedTemplates[name]; ok {
		return t
	}

	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles(
		filepath.Join(Settings.Root, "templates/_base.html"),
		filepath.Join(Settings.Root, "templates", name),
	))
	cachedTemplates[name] = t

	return t
}

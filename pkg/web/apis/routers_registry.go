package apis

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// Strapper ...
type Strapper interface {
	Strap(r gin.IRouter)
}

// StrapFunc ...
type StrapFunc func(r gin.IRouter)

// Strap ...
func (f StrapFunc) Strap(r gin.IRouter) {
	f(r)
}

var straps = make(map[string]Strapper)

// Register ...
func Register(name string, sf Strapper) {
	straps[name] = sf
}

// Routers ...
func Routers(r gin.IRouter, names ...string) {
	logger().Infow("Routers", "names", names)
	var once sync.Once
	once.Do(func() {
		if len(names) == 0 {
			for name, sf := range straps {
				logger().Infow("start router", "name", name)
				sf.Strap(r)
			}
			return
		}
		for _, name := range names {
			if sf, ok := straps[name]; ok {
				logger().Infow("start router", "name", name)
				sf.Strap(r)
			} else {
				logger().Warnw("strap not found", "name", name)
			}
		}
	})
}

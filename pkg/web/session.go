package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-osin/session"
	scodec "github.com/go-osin/session/codec"
	"github.com/go-osin/session/redicache"
	"github.com/ugorji/go/codec"

	"github.com/liut/staffio/pkg/settings"
)

func init() {
	if len(settings.Current.RedisAddrs) > 0 {
		setupSessionStore(redicache.NewStoreOptions(&redicache.StoreOptions{
			Codec:     &MsgPack,
			Addrs:     settings.Current.RedisAddrs,
			DB:        settings.Current.RedisDB,
			Password:  settings.Current.RedisPassword,
			KeyPrefix: "staffio-sess-",
		}))
	}

}

var (
	sessionKey = "gin-session"
)

var (
	// MsgPack is a Codec that uses the `ugorji/go/codec` package.
	MsgPack = scodec.Codec{Marshal: msgPackMarshal, Unmarshal: msgPackUnmarshal}
)

func msgPackMarshal(v interface{}) (out []byte, err error) {
	var h codec.Handle = new(codec.MsgpackHandle)
	err = codec.NewEncoderBytes(&out, h).Encode(v)
	return
}

func msgPackUnmarshal(in []byte, v interface{}) error {
	var h codec.Handle = new(codec.MsgpackHandle)
	return codec.NewDecoderBytes(in, h).Decode(v)
}

func setupSessionStore(store session.Store) {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(store, &session.CookieMngrOptions{
		SessIDCookieName: "st_sess",
		AllowHTTP:        true,
	})
}

func SessionLoad(r *http.Request) session.Session {
	sess := session.Global.Load(r)
	if sess == nil {
		sess = session.NewSession()
	}
	return sess
}

func SessionSave(sess session.Session, w http.ResponseWriter) {
	session.Global.Save(sess, w)
}

func ginSession(c *gin.Context) session.Session {
	if sess, ok := c.Get(sessionKey); ok {
		return sess.(session.Session)
	}
	sess := SessionLoad(c.Request)
	c.Set(sessionKey, sess)
	return sess
}

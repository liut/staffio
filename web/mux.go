package web

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/getsentry/raven-go"
	"github.com/goods/httpbuf"

	. "lcgc/platform/staffio/settings"
)

type handler func(*Context) error

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Expires", "Fri, 02 Oct 1998 20:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate")

	origin := req.Header.Get("Origin")
	if len(origin) > 0 {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "60")

		if req.Method == "OPTIONS" {
			w.Header().Set("Allow", "GET,HEAD,POST,OPTIONS")
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	buf := new(httpbuf.Buffer)

	sess, err := store.Get(req, Settings.Session.Name)
	if err != nil {
		// debug.PrintStack()
		raven.CaptureError(err, nil)
		log.Printf("load session error: %s", err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
	}
	ctx := NewContext(buf, req, sess)
	defer ctx.Close()

	//run the handler and grab the error, and report it
	err = h(ctx)
	if err != nil {
		debug.PrintStack()
		if ctx.User != nil {
			raven.SetUserContext(&raven.User{ID: ctx.User.Uid})
		}
		raven.SetHttpContext(raven.NewHttp(req))
		logId := raven.CaptureError(err, nil)
		raven.ClearContext()
		log.Printf("[%s] %s %s error: %s logId: %s", req.RemoteAddr, req.Method, req.RequestURI, err, logId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if req.Method != http.MethodGet {
			log.Printf("[%s] %s %s OK", req.RemoteAddr, req.Method, req.RequestURI)
		}
	}

	ctx.afterHandle()

	//save the session
	if len(ctx.Session.Values) > 0 { // session not empty only
		if err = ctx.Session.Save(req, buf); err != nil {
			log.Printf("session.save error: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//apply the buffered response to the writer
	buf.Apply(w)
}

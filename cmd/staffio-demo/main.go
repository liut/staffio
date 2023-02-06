package main

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	talog "daxv.cn/gopak/tencent-api-go/log"
	"github.com/liut/keeper/utils/reaper"
	"github.com/liut/staffio/pkg/backends"
	zlog "github.com/liut/staffio/pkg/log"
	config "github.com/liut/staffio/pkg/settings"
	"github.com/liut/staffio/pkg/web"
)

const (
	readTimeout  time.Duration = 10 * time.Second
	writeTimeout               = 15 * time.Second
)

func main() {
	var zlogger *zap.Logger

	zlogger, _ = zap.NewDevelopment()

	defer zlogger.Sync() // flushes buffer, if any
	sugar := zlogger.Sugar()

	zlog.SetLogger(sugar)
	talog.SetLogger(sugar)

	log.SetFlags(log.Ltime | log.Lshortfile)

	settings := config.Current

	cfg := web.Config{
		Root:    ".",
		FS:      "local",
		BaseURI: settings.BaseURL,
	}
	ws := web.New(cfg)
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	srv := &http.Server{
		Addr:         settings.HTTPListen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      ws,
	}

	d := &demo{
		base: settings.BaseURL,
	}
	d.strap(ws)

	logger().Infow("Starting", "ver", config.Version(), "listen", settings.HTTPListen, "root", settings.Root)
	err := srv.ListenAndServe() // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}
}

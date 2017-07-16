//
// Main process for run web server
//
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wealthworks/csmtp"
	"github.com/wealthworks/go-utils/reaper"

	"lcgc/platform/staffio/pkg/backends"
	. "lcgc/platform/staffio/pkg/settings"
	"lcgc/platform/staffio/pkg/web"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	Settings.Parse()

	csmtp.Host = Settings.SMTP.Host
	csmtp.Port = Settings.SMTP.Port
	csmtp.Name = Settings.SMTP.SenderName
	csmtp.From = Settings.SMTP.SenderEmail
	csmtp.Auth(Settings.SMTP.SenderPassword)

	ws := web.New()
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)

	srv := &http.Server{
		Addr:    Settings.HttpListen,
		Handler: ws,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Print("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Print("Server exit")
}

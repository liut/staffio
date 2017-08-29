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

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
	"github.com/liut/staffio/pkg/web"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	settings.Parse()

	csmtp.Host = settings.SMTP.Host
	csmtp.Port = settings.SMTP.Port
	csmtp.Name = settings.SMTP.SenderName
	csmtp.From = settings.SMTP.SenderEmail
	csmtp.Auth(settings.SMTP.SenderPassword)

	ws := web.New()
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", settings.Version(), settings.HttpListen, settings.Root)

	srv := &http.Server{
		Addr:    settings.HttpListen,
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

// Copyright © 2019 liut <liutao@liut.cc>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/liut/keeper/utils/reaper"
	"github.com/spf13/cobra"

	"github.com/liut/staffio/pkg/backends"
	config "github.com/liut/staffio/pkg/settings"
	"github.com/liut/staffio/pkg/web"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start main web server",
	Long:  `Start main web server`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.ParseFlags(args)
		webRun()
	},
}

var webFS string
var webRoot string

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVar(&webFS, "fs", "bind", "file system [bind | local]")
	webCmd.Flags().StringVar(&webRoot, "root", "./", "app root directory")
}

const (
	readTimeout  time.Duration = 10 * time.Second
	writeTimeout               = 15 * time.Second
)

func webRun() {
	// web.SetBase("/v1/")
	cfg := web.Config{
		Root:    webRoot,
		FS:      webFS,
		BaseURI: settings.BaseURL,
	}
	ws := web.New(cfg)
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", config.Version(), settings.HTTPListen, settings.Root)

	srv := &http.Server{
		Addr:         settings.HTTPListen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      ws,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen failed: %s", err)
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

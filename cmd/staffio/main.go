//
// Main process for run web server
//
package main

import (
	"fmt"
	"log"

	"github.com/wealthworks/go-utils/reaper"

	"lcgc/platform/staffio/pkg/backends"
	. "lcgc/platform/staffio/pkg/settings"
	"lcgc/platform/staffio/pkg/web"
)

func main() {
	ws := web.New()
	defer reaper.Quit(reaper.Run(0, backends.Cleanup))

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)
	err := ws.Run(Settings.HttpListen) // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}

}

package main

import (
	"fmt"
	"log"
	"strings"

	. "lcgc/platform/staffio/pkg/settings"
	"lcgc/platform/staffio/pkg/web"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	Settings.Parse()
	ws := web.New()
	if strings.HasPrefix(Settings.HttpListen, "localhost") {
		d := &demo{
			prefix: "http://" + Settings.HttpListen,
		}
		d.strap(ws)
	}

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)
	err := ws.Run(Settings.HttpListen) // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}
}

package main

import (
	"fmt"
	"log"
	"strings"

	"lcgc/platform/staffio/pkg/settings"
	"lcgc/platform/staffio/pkg/web"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	settings.Parse()
	ws := web.New()
	if strings.HasPrefix(settings.HttpListen, "localhost") {
		d := &demo{
			prefix: "http://" + settings.HttpListen,
		}
		d.strap(ws)
	}

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", settings.Version, settings.HttpListen, settings.Root)
	err := ws.Run(settings.HttpListen) // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}
}

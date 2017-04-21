package main

import (
	"fmt"
	"log"
	"strings"

	. "lcgc/platform/staffio/settings"
	"lcgc/platform/staffio/web"
)

func main() {
	ws := web.New()
	if strings.HasPrefix(Settings.HttpListen, "localhost") {
		AppDemo(ws)
	}

	fmt.Printf("Start service %s at addr %s\nRoot: %s\n", Settings.Version, Settings.HttpListen, Settings.Root)
	err := ws.Run(Settings.HttpListen) // Start the server!
	if err != nil {
		log.Fatal("Run ERR: ", err)
	}
}

package main

import (
	"flag"
	"log"

	"lcgc/platform/staffio/backends/exmail"
	. "lcgc/platform/staffio/settings"
)

var (
	alias = flag.String("alias", "", "alias: username | email")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Settings.Parse()
	if *alias == "" {
		flag.PrintDefaults()
		return
	}

	staff, err := exmail.GetStaff(*alias)
	if err != nil {
		log.Printf("get staff err: %s", err)
		return
	}
	log.Printf("staff: %v", staff)
}

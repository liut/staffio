package main

import (
	"flag"
	"log"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
	"github.com/wealthworks/go-tencent-api/exmail"
)

var (
	alias = flag.String("alias", "", "alias: username | email")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	settings.Parse()
	if *alias == "" {
		flag.PrintDefaults()
		return
	}

	staff, err := backends.GetStaffFromExmail(*alias)
	if err != nil {
		log.Printf("get staff err: %s", err)
		return
	}
	log.Printf("staff: %v", staff)

	count, err := exmail.CountNewMail(*alias)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("new mail: %d", count)

	url, err := exmail.GetLoginURL(*alias)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("open %s", url)
}

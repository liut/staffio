package main

import (
	"flag"
	"log"

	"github.com/wealthworks/go-tencent-api/exmail"
	"github.com/wealthworks/go-tencent-api/exwechat"
	"lcgc/platform/staffio/backends"
	. "lcgc/platform/staffio/settings"
)

var (
	uid = flag.String("uid", "", "uid: username | email")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Settings.Parse()
	backends.Prepare()
	if *uid == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("fetching user %s", *uid)

	staff, err := backends.GetStaff(*uid)
	if err != nil {
		log.Printf("get staff err: %s", err)
		return
	}
	log.Printf("staff: %v", staff)
	if staff == nil {
		return
	}

	wechat := exwechat.NewAPI()
	user, err := wechat.GetUser(*uid)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("WxWork User: %v", user)

	count, err := exmail.CountNewMail(*uid)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("new mail: %d", count)

	url, err := exmail.GetLoginURL(*uid)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("open %s", url)
}

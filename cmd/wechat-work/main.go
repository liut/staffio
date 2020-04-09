// List or sync teams data from department of wechat work
package main

import (
	"flag"
	"log"

	"github.com/liut/staffio/pkg/backends/wechatwork"
)

var (
	action string
	uid    string
)

func init() {
	flag.StringVar(&action, "act", "", "action: list | query | sync | sync-all")
	flag.StringVar(&uid, "uid", "", "query uid")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	if action == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("action: %q", action)

	// backends.InitSMTP()

	wechatwork.SyncDepartment(action, uid)
}

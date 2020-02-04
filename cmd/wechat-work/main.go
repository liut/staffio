// List or sync teams data from department of wechat work
package main

import (
	"flag"
	// "fmt"
	"log"
	"strings"
	// "time"

	// "github.com/wealthworks/go-tencent-api/exwechat"

	// "github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/backends/wechatwork"
	// "github.com/liut/staffio/pkg/models"
	// "github.com/liut/staffio/pkg/settings"
)

var (
	action string
	uid    string

	nameReplacer = strings.NewReplacer("公司", "", "总部", "")
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

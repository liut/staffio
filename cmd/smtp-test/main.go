package main

import (
	"flag"
	"log"

	"github.com/wealthworks/csmtp"
	. "lcgc/platform/staffio/pkg/settings"
)

var (
	toEmail string
)

func init() {
	flag.StringVar(&toEmail, "to", "", "")
}

func main() {
	Settings.Parse()
	if toEmail == "" {
		flag.PrintDefaults()
		return
	}

	csmtp.Host = Settings.SMTP.Host
	csmtp.Port = Settings.SMTP.Port
	csmtp.Name = Settings.SMTP.SenderName
	csmtp.From = Settings.SMTP.SenderEmail
	csmtp.Auth(Settings.SMTP.SenderPassword)

	subject := "测试主题"
	body := "我是一封电子邮件!golang发出."

	err := csmtp.SendMail(subject, body, toEmail)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("send OK")
	}
}

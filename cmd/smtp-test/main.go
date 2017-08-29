package main

import (
	"flag"
	"log"

	"github.com/liut/staffio/pkg/settings"
	"github.com/wealthworks/csmtp"
)

var (
	toEmail string
)

func init() {
	flag.StringVar(&toEmail, "to", "", "")
}

func main() {
	settings.Parse()
	if toEmail == "" {
		flag.PrintDefaults()
		return
	}

	csmtp.Host = settings.SMTP.Host
	csmtp.Port = settings.SMTP.Port
	csmtp.Name = settings.SMTP.SenderName
	csmtp.From = settings.SMTP.SenderEmail
	csmtp.Auth(settings.SMTP.SenderPassword)

	subject := "测试主题"
	body := "我是一封电子邮件!golang发出."

	err := csmtp.SendMail(subject, body, toEmail)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("send OK")
	}
}

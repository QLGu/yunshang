package mail

import (
	"log"
	"strings"

	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

var _rules = map[string]string{"gmail.com": "mail.google.com", "139.com": "mail.10086.cn"}

func ModuleInit() {
	log.Printf("Init Module %v", "mail")

	provider := revel.Config.StringDefault("mail.provider", "google")
	switch provider {
	case "google":
		email.InitGmail("yunshang2014@gmail.com", "revel2014")
	case "qq":
		initQQMailFrom("cljwtang@qq.com", "cljwtang@2013")
	}
	email.Config.From.Name = revel.Config.StringDefault("mail.fromName", "YuShang")
}

func Send(subject, content, to string) error {
	return send(subject, content, to, false)
}

func SendHtml(subject, content, to string) error {
	return send(subject, content, to, true)
}

func GetEmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	host := arrs[1]
	provider, ok := _rules[host]
	if ok {
		return "http://" + provider
	}
	return "http://mail." + host
}

func send(subject, content, to string, html bool) error {
	mail := email.NewBriefMessage(subject, content, to)
	mail.IsHtmlContent = html
	return mail.Send()
}

func initQQMailFrom(fromAddress, password string) (err error) {
	if err = email.InitGmailFrom(fromAddress, fromAddress, password); err != nil {
		return
	}
	email.Config.Host = "smtp.qq.com"
	email.Config.Port = 465
	return
}

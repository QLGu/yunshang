package controllers

import (
	"fmt"
	"strings"

	"github.com/ungerik/go-mail"
)

func SendMail(subject, content, to string) {
	sendMail(subject, content, to, false)
}

func SendHtmlMail(subject, content, to string) {
	sendMail(subject, content, to, true)
}

func sendMail(subject, content, to string, html bool) {
	mail := email.NewBriefMessage(subject, content, to)
	mail.IsHtmlContent = html
	err := mail.Send()

	if err != nil {
		fmt.Println(err)
	}
}

func EmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	rules := map[string]string{"gmail.com":"mail.google.com", "139.com": "mail.10086.cn"}
	host := arrs[1]
	provider, ok := rules[host]
	if ok {
		return "http://" + provider
	}
	return "http://mail." + host
}

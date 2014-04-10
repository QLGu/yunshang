//
// 邮件处理模块.
//
package models

import (
	"log"
	"strconv"
	"strings"

	"github.com/itang/gotang"
	"github.com/ungerik/go-mail"
)

// 邮箱后缀对应的服务地址
var _mail_rules = map[string]string{"gmail.com": "mail.google.com", "139.com": "mail.10086.cn"}

func InitMailConfig() {
	log.Println("InitMailConfig")

	email.Config.Host = CacheSystem.GetConfig("site.mail.host")

	port, err := strconv.Atoi(CacheSystem.GetConfig("site.mail.port"))
	gotang.AssertNoError(err, "site.mail.port号设置不正确")
	email.Config.Port = uint16(port)

	email.Config.Username = CacheSystem.GetConfig("site.mail.username")
	email.Config.Password = CacheSystem.GetConfig("site.mail.password")

	email.Config.From.Name = CacheSystem.GetConfig("site.mail.from_name")
	email.Config.From.Address = CacheSystem.GetConfig("site.mail.from_address")
}

// 发送邮件
func SendMail(subject, content, to string) error {
	return sendMail(subject, content, to, false)
}

// 发送HTML邮件
func SendHtmlMail(subject, content, to string) error {
	return sendMail(subject, content, to, true)
}

// 获取邮箱对应服务的WEB URL
func GetEmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	host := arrs[1]
	provider, ok := _mail_rules[host]
	if ok {
		return "http://" + provider
	}
	return "http://mail." + host
}

// 发送邮件
func sendMail(subject, content, to string, html bool) error {
	mail := email.NewBriefMessage(subject, content, to)
	mail.IsHtmlContent = html
	return mail.Send()
}

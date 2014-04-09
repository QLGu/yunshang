//
// 邮件处理模块.
//
package mail

import (
	"log"
	"strings"

	"github.com/itang/yunshang/main/app/cache"
	"github.com/revel/revel"
	"github.com/ungerik/go-mail"
)

// 邮箱后缀对应的服务地址
var _rules = map[string]string{"gmail.com": "mail.google.com", "139.com": "mail.10086.cn"}

// 模块注册
func ModuleInit() {
	log.Printf("Init Module %v", "mail")

	provider := revel.Config.StringDefault("mail.provider", "google")
	switch provider {
	case "google":
		email.InitGmail("yunshang2014@gmail.com", "revel2014")
	case "qq":
		initQQMailFrom("cljwtang@qq.com", "cljwtang@2013")
	}
	email.Config.From.Name = cache.GetConfig("site.basic.name") //revel.Config.StringDefault("mail.fromName", "YuShang")
}

// 发送邮件
func Send(subject, content, to string) error {
	return send(subject, content, to, false)
}

// 发送HTML邮件
func SendHtml(subject, content, to string) error {
	return send(subject, content, to, true)
}

// 获取邮箱对应服务的WEB URL
func GetEmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	host := arrs[1]
	provider, ok := _rules[host]
	if ok {
		return "http://" + provider
	}
	return "http://mail." + host
}

// 发送邮件
func send(subject, content, to string, html bool) error {
	mail := email.NewBriefMessage(subject, content, to)
	mail.IsHtmlContent = html
	return mail.Send()
}

// 初始化QQ-SMTP
func initQQMailFrom(fromAddress, password string) (err error) {
	if err = email.InitGmailFrom(fromAddress, fromAddress, password); err != nil {
		return
	}
	email.Config.Host = "smtp.qq.com"
	email.Config.Port = 465
	return
}

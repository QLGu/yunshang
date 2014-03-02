package controllers

import (
	"strconv"
	"strings"

	"github.com/ungerik/go-mail"
)

func SendMail(subject, content, to string) error {
	return sendMail(subject, content, to, false)
}

func SendHtmlMail(subject, content, to string) error {
	return sendMail(subject, content, to, true)
}

func sendMail(subject, content, to string, html bool) error {
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

func EmailProvider(email string) string {
	arrs := strings.Split(email, "@")
	rules := map[string]string{"gmail.com": "mail.google.com", "139.com": "mail.10086.cn"}
	host := arrs[1]
	provider, ok := rules[host]
	if ok {
		return "http://" + provider
	}
	return "http://mail." + host
}

type dataTableData struct {
	SEcho                int         `json:"sEcho"`
	ITotalRecords        int64       `json:"iTotalRecords"`
	ITotalDisplayRecords int64       `json:"iTotalDisplayRecords"`
	AaData               interface{} `json:"aaData,omitempty"`
}

func DataTableData(echo string, total int64, totalDisplay int64, data interface{}) dataTableData {
	ei, err := strconv.Atoi(echo)
	if err != nil {
		ei = 0
	}
	return dataTableData{SEcho: ei, ITotalRecords: total, ITotalDisplayRecords: totalDisplay, AaData: data}
}

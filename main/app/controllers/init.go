package controllers

import (
	"code.google.com/p/goauth2/oauth"

	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

var QQ = &oauth.Config{
	ClientId:     "1101230147",
	ClientSecret: "nUAwR4KCygCsEtNM",
	AuthURL:      "https://graph.facebook.com/oauth/authorize",
	TokenURL:     "https://graph.facebook.com/oauth/access_token",
	RedirectURL:  "http://loisant.org:9000/Application/Auth",
}

func init() {
	revel.InterceptMethod((*XOrmController).begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).rollback, revel.PANIC)

	revel.InterceptMethod((*ShouldLoginedController).checkUser, revel.BEFORE)
	revel.InterceptMethod((*AdminController).checkAdminUser, revel.BEFORE)

	email.InitGmail("yunshang2014@gmail.com", "revel2014")
	email.Config.From.Name = "YuShang"
}

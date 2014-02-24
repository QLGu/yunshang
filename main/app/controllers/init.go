package controllers

import (
	"code.google.com/p/goauth2/oauth"

	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

var WEIBO = &oauth.Config{
	ClientId:       "3635003648",
	ClientSecret:   "76d75007b6ea3b05e24762105e773270",
	Scope:          "email",
	AccessType:     "offline",
	ApprovalPrompt: "auto",
	AuthURL:        "https://api.weibo.com/oauth2/authorize",
	TokenURL:       "https://api.weibo.com/oauth2/access_token",
	RedirectURL:    "http://yunshang.haoshuju.net/passport/open/weibo",
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

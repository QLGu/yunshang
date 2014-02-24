package controllers

import (
	"fmt"
	"strconv"

	"github.com/itang/yunshang/modules/oauth"
	"github.com/itang/yunshang/modules/oauth/apps"

	"github.com/itang/gotang"
	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

var WEIBO = &oauth.Config{
	ClientId:     "3635003648",
	ClientSecret: "76d75007b6ea3b05e24762105e773270",
	//Scope:          "email",
	//AccessType:     "offline",
	//ApprovalPrompt: "auto",
	AuthURL:     "https://api.weibo.com/oauth2/authorize",
	TokenURL:    "https://api.weibo.com/oauth2/access_token",
	RedirectURL: "http://yunshang.haoshuju.net/passport/open/weibo",
}

var SocialAuth *oauth.SocialAuth

func init() {
	revel.InterceptMethod((*XOrmController).begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).rollback, revel.PANIC)

	revel.InterceptMethod((*ShouldLoginedController).checkUser, revel.BEFORE)
	revel.InterceptMethod((*AdminController).checkAdminUser, revel.BEFORE)

	email.InitGmail("yunshang2014@gmail.com", "revel2014")
	email.Config.From.Name = "YuShang"

	revel.OnAppStart(initOAuth)
}

func initOAuth() {
	// OAuth
	var clientId, secret string

	var err error
	appURL := revel.Config.StringDefault("social_auth_url", "")
	if len(appURL) > 0 {
		oauth.DefaultAppUrl = appURL
	}

	clientId = revel.Config.StringDefault("weibo_client_id", "")
	secret = revel.Config.StringDefault("weibo_client_secret", "")
	err = oauth.RegisterProvider(apps.NewWeibo(clientId, secret))
	gotang.AssertNoError(err)

	//clientId = revel.Config.StringDefault("qq_client_id","")
	//secret = revel.Config.StringDefault("qq_client_secret","")
	//err = oauth.RegisterProvider(apps.NewQQ(clientId, secret))
	gotang.AssertNoError(err)

	SocialAuth = oauth.NewSocial("/passport/open/", new(socialAuther))
}

type socialAuther struct {
}

func (p *socialAuther) IsUserLogin(ctx *revel.Controller) (int, bool) {
	us, ok := ctx.Session["user"]
	i, err := strconv.Atoi(us)
	return i, ok && err == nil
}

func (p *socialAuther) LoginUser(ctx *revel.Controller, uid int) (string, error) {
	ctx.Session["user"] = fmt.Sprintf("%v", uid)

	return "/", nil
}

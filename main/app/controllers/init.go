package controllers

import (
	//"fmt"
	"strconv"

	"github.com/itang/gotang"
	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"

	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/modules/oauth"
	"github.com/itang/yunshang/modules/oauth/apps"
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
	//initQQMailFrom("cljwtang@qq.com", "cljwtang@2013")
	email.Config.From.Name = "YuShang"

	revel.OnAppStart(initOAuth)
}

func initOAuth() {
	// OAuth
	var clientId, secret string

	appURL := revel.Config.StringDefault("social_auth_url", "http://"+revel.Config.StringDefault("web.host", ""))
	if len(appURL) > 0 {
		oauth.DefaultAppUrl = appURL
	}

	clientId = revel.Config.StringDefault("weibo_client_id", "")
	secret = revel.Config.StringDefault("weibo_client_secret", "")
	gotang.Assert(clientId != "" && secret != "", "weibo_client_id和weibo_client_secret不能为空")

	err := oauth.RegisterProvider(apps.NewWeibo(clientId, secret))
	gotang.AssertNoError(err)

	//clientId = revel.Config.StringDefault("qq_client_id","")
	//secret = revel.Config.StringDefault("qq_client_secret","")
	//err = oauth.RegisterProvider(apps.NewQQ(clientId, secret))

	SocialAuth = oauth.NewSocial("/passport/open/", new(socialAuther))
}

type socialAuther struct {
}

func (p *socialAuther) IsUserLogin(ctx *revel.Controller) (int64, bool) {
	us, ok := ctx.Session["uid"]

	i, err := strconv.Atoi(us)
	return int64(i), ok && err == nil
}

func (p *socialAuther) LoginUser(ctx *revel.Controller, uid int64, socialType oauth.SocialType) (string, error) {
	passport, ok := ctx.AppController.(*Passport)
	gotang.Assert(ok, "FROM passport")

	user, ok := passport.userService().GetUserById(uid)
	revel.INFO.Printf("user:id %v, %v", user.Id, user)
	gotang.Assert(ok, "user not exists")
	if !ok || !user.Enabled {
		passport.Flash.Error("用户信息不存在或被禁用，有任何疑问请联系本站客服！")
		return "/passport/login", nil
	}
	passport.SetLoginSession(models.ToSessionUser(user))
	// 执行登录后操作
	go passport.userService().DoUserLogin(&user)

	return "/", nil
}

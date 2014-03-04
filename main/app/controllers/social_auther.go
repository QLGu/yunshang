package controllers

import (
	"strconv"

	"github.com/itang/gotang"
	"github.com/revel/revel"

	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/modules/oauth"
)

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
	passport.setLoginSession(models.ToSessionUser(user))
	// 执行登录后操作
	go passport.userService().DoUserLogin(&user)

	return "/", nil
}

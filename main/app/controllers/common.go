package controllers

import (
	"fmt"
	"strings"

	"github.com/robfig/revel"

	"github.com/itang/gotang"
	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"strconv"
)

type RestResposne struct {
	Ok      bool        `json:"ok"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AppController struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
}

func (c AppController) successResposne(message string, data interface{}) RestResposne {
	return RestResposne{Ok: true, Code: 0, Message: message, Data: data}
}

func (c AppController) errorResposne(message string, data interface{}) RestResposne {
	return RestResposne{Ok: false, Code: 1, Message: message, Data: data}
}

func (c AppController) isLogined() bool {
	_, ok := c.Session["uid"]
	return ok
}

func (c AppController) currUser() (user entity.User, ok bool) {
	uidStr, ok := c.Session["uid"]
	if !ok {
		return user, false
	}

	id, err := strconv.Atoi(uidStr)
	if err != nil {
		return user, false
	}
	return c.userService().GetUserById(int64(id))
}

func (c AppController) SetLoginSession(sessionUser models.SessionUser) {
	c.Session["login"] = sessionUser.LoginName
	c.Session["uid"] = fmt.Sprintf("%v", sessionUser.Id)
	c.Session["screen_name"] = sessionUser.DisplayName()
	c.Session["from"] = sessionUser.From
}

func (c AppController) ClearLoginSession() {
	delete(c.Session, "login")
	delete(c.Session, "uid")
	delete(c.Session, "screen_name")
	delete(c.Session, "from")
}

func (c AppController) forceCurrUser() (user entity.User) {
	user, ok := c.currUser()
	gotang.Assert(ok, "用户未登录！")
	return
}

func (c AppController) doValidate(redirectTarget interface{}) revel.Result {
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(redirectTarget)
	}
	return nil
}

func (c AppController) userService() models.UserService {
	gotang.Assert(c.XOrmSession != nil, "c.XOrmSession should no be nil")
	return models.DefaultUserService(c.XOrmSession)
}

func (c AppController) getRemoteIp() string {
	ips, ok := c.Request.Header["X-Real-IP"]
	if !ok {
		return strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	return ips[0]
}

type ShouldLoginedController struct {
	AppController
}

func (c ShouldLoginedController) checkUser() revel.Result {
	if !c.isLogined() {
		return c.Redirect(App.Index)
	}
	return nil
}

type AdminController struct {
	ShouldLoginedController
}

func (c AdminController) checkAdminUser() revel.Result {
	user, _ := c.Session["screen_name"]
	//TODO 使用角色
	if !c.isLogined() || user != "admin" {
		return c.Redirect(App.Index)
	}
	return nil
}

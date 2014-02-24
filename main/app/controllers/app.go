package controllers

import (
	"strings"

	"github.com/robfig/revel"

	"github.com/itang/iptaobao"
	"github.com/itang/gotang"
	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
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
	_, ok := c.Session["user"]
	return ok
}

func (c AppController) currUser() (user entity.User, ok bool) {
	login, ok := c.Session["user"]
	if ok {
		return c.userService().GetUserByLogin(login)
	} else {
		return user, false
	}
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
	user, _ := c.Session["user"]
	//TODO 使用角色
	if !c.isLogined() || user != "admin" {
		return c.Redirect(App.Index)
	}
	return nil
}

type App struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
	AppController
}

func (c App) Index() revel.Result {
	c.RenderArgs["version"] = app.Version


	ip := strings.Split(c.Request.RemoteAddr, ":")[0]
	revel.INFO.Printf("remoteAddr: %v", ip)

	c.RenderArgs["ip"] = ip

	r, err := iptaobao.GetIpInfo(ip)
	if err != nil {
		revel.INFO.Printf("%v, %v", ip, err)
	}else {
		c.RenderArgs["from"] = r.City + r.Area + r.Region
	}

	return c.Render()
}

func (c App) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c App) Panic() revel.Result {
	return c.RenderText("hello")
}

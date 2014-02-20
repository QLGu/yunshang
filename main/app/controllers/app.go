package controllers

import (
	"github.com/robfig/revel"

	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models"
)

type AppController struct {
	*revel.Controller
}

func (c AppController) IsLogined() bool {
	_, ok := c.Session["user"]
	return ok
}

type App struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
	AppController
}

func (c App) Index() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()
	c.RenderArgs["version"] = app.Version

	return c.Render()
}

func (c App) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c App) Panic() revel.Result {
	return c.RenderText("hello")
}

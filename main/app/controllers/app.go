package controllers

import (
	"github.com/robfig/revel"

	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models"
)

type App struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
}

func (c App) Index() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()
	c.RenderArgs["version"] = app.Version

	//test
	//c.Session["user"]="itang"

	return c.Render()
}

func (c App) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

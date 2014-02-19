package controllers

import (
	"github.com/robfig/revel"

	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/models"
)

type Admin struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
}

func (c Admin) Users() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()

	return c.Render()
}

func (c Admin) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

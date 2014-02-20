package controllers

import (
	"github.com/robfig/revel"

	//reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/utils"
)

type Admin struct {
	AdminController
}

func (c Admin) Users() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()

	return c.Render()
}

type JsonResult struct {
	Ok      bool
	Message string
	Data    interface{}
}

func (c Admin) ResetUserPassword(id int64) revel.Result {
	user, ok := c.userService().GetUserById(id)
	if !ok {
		return c.RenderJson(JsonResult{false, "用户不存在", nil})
	}
	newPassword := utils.RandomString(6)
	err := c.userService().DoChangePassword(&user, newPassword)
	if err != nil {
		return c.RenderJson(JsonResult{false, err.Error(), nil})
	}

	data := struct {
		NewPassword string
	}{newPassword}
	/*go*/ SendHtmlMail("重置密码邮件", utils.RenderTemplateToString("Passport/ResetPasswordResultTemplate.html", data), user.Email)

	return c.RenderJson(JsonResult{true, "重置用户密码成功！", newPassword})
}

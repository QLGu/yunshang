package controllers

import (
	"github.com/robfig/revel"

	//reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/utils"
)

type Admin struct {
	AdminController
}

func (c Admin) Index() revel.Result {
	_, ok := c.Session["locked"]
	if ok {
		c.Redirect(Admin.Lock)
	}

	c.RenderArgs["channel"] = "/"
	return c.Render()
}

func (c Admin) Lock() revel.Result {
	c.Session["locked"] = "true"
	return c.Render()
}

func (c Admin) UnLock(password string) revel.Result {
	delete(c.Session, "locked")
	return c.Redirect(Admin.Index)
}

func (c Admin) Users() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()
	c.RenderArgs["channel"] = "users/users"

	return c.Render()
}

func (c Admin) UsersData() revel.Result {
	page := c.userService().FindAllUsersForPage(c.pageSearcher())
	return c.renderDataTableJson(page)
}

func (c Admin) ResetUserPassword(id int64) revel.Result {
	user, ok := c.userService().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}
	newPassword := utils.RandomString(6)
	err := c.userService().DoChangePassword(&user, newPassword)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	}

	data := struct {
		NewPassword string
	}{newPassword}

	/*go*/ SendHtmlMail("重置密码邮件", utils.RenderTemplateToString("Passport/ResetPasswordResultTemplate.html", data), user.Email)

	return c.RenderJson(c.successResposne("重置用户密码成功并新密码已经通过告知邮件用户", newPassword))
}

func (c Admin) ToggleUserEnabled(id int64) revel.Result {
	user, ok := c.userService().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}

	if c.userService().IsAdminUser(&user) {
		return c.RenderJson(c.errorResposne("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userService().ToggleUserEnabled(&user)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		return c.RenderJson(c.successResposne("改变用户状态！", nil))
	}
}

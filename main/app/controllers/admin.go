package controllers

import (
	"time"

	"github.com/robfig/revel"

	//reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/mail"
)

// 管理端相关Actions
type Admin struct {
	AdminController
}

// 管理端主页
func (c Admin) Index() revel.Result {
	_, ok := c.Session["locked"]
	if ok {
		c.Redirect(Admin.Lock)
	}

	c.RenderArgs["channel"] = "/"
	return c.Render()
}

// 锁屏
func (c Admin) Lock() revel.Result {
	c.Session["locked"] = "true"
	return c.Render()
}

// 解锁屏
func (c Admin) UnLock(password string) revel.Result {
	delete(c.Session, "locked")
	return c.Redirect(Admin.Index)
}

// 用户列表
func (c Admin) Users() revel.Result {
	c.RenderArgs["users_total"] = c.userService().Total()
	c.RenderArgs["users"] = c.userService().FindAllUsers()
	c.RenderArgs["channel"] = "users/users"

	return c.Render()
}

// 用户列表数据
func (c Admin) UsersData() revel.Result {
	page := c.userService().FindAllUsersForPage(c.pageSearcher())
	return c.renderDataTableJson(page)
}

// 重置用户密码
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

	err = utils.DoIOWithTimeout(func() error {
		return mail.SendHtml("重置密码邮件",
			utils.RenderTemplate("Passport/ResetPasswordResultTemplate.html",
				struct {
					NewPassword string
				}{newPassword}),
			user.Email)
	}, time.Second*30)
	if err != nil {
		panic(err)
	}

	return c.RenderJson(c.successResposne("重置用户密码成功并新密码已经通过告知邮件用户", newPassword))
}

// 激活/禁用用户
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

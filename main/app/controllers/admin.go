package controllers

import (
	"time"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/mail"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
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

	userTotal := c.userService().Total()
	orderTotal := 0 // TODO order total

	c.setChannel("/")
	return c.Render(userTotal, orderTotal)
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
	c.setChannel("users/users")
	return c.Render()
}

// 用户列表数据
func (c Admin) UsersData(filter_status string, filter_certified string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}

		switch filter_certified {
		case "true":
			session.And("certified=?", true)
		case "false":
			session.And("certified=?", false)
		}
	})
	page := c.userService().FindAllUsersForPage(ps)
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

	err = gotang.DoIOWithTimeout(func() error {
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

// 认证用户
func (c Admin) ToggleUserCertified(id int64) revel.Result {
	user, ok := c.userService().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}

	if c.userService().IsAdminUser(&user) {
		return c.RenderJson(c.errorResposne("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userService().ToggleUserCertified(&user)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		return c.RenderJson(c.successResposne("改变用户状态！", nil))
	}
}

// 显示用户登录日志
func (c Admin) ShowUserLoginLogs(id int64) revel.Result {
	loginLogs := c.userService().FindUserLoginLogs(id)
	return c.Render(loginLogs)
}

// 显示用户信息
func (c Admin) ShowUserInfos(id int64) revel.Result {
	user, _ := c.userService().GetUserById(id)
	userDetail, _ := c.userService().GetUserDetailByUserId(user.Id)
	userDas := c.userService().FindUserDeliveryAddresses(user.Id)

	return c.Render(user, userDetail, userDas)
}

// 产品列表
func (c Admin) Products() revel.Result {
	c.setChannel("products/products")
	return c.Render()
}

// 产品列表数据
func (c Admin) ProductsData(filter_status string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
	})
	page := c.productApi().FindAllProductsForPage(ps)
	return c.renderDataTableJson(page)
}

func (c Admin) ToggleProductEnabled() revel.Result {
	return c.RenderJson("")
}

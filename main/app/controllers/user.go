package controllers

import (
	"strings"

	"github.com/robfig/revel"

	//reveltang "github.com/itang/reveltang/controllers"
)

type User struct {
	ShouldLoginedController
}

func (c User) Index() revel.Result {
	currUser, _:= c.currUser()
	return c.Render(currUser)
}

func (c User) ChangePassword() revel.Result {
	return c.Render()
}

func (c User) DoChangePassword(oldPassword, password, confirmPassword string) revel.Result {
	c.Validation.Required(oldPassword).Message("请输入密码")
	c.Validation.Required(password).Message("请输入新密码")
	c.Validation.MinSize(password, 6).Message("请输入6位新密码")
	c.Validation.MaxSize(password, 50).Message("输入新密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认新密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的新密码不匹配").Key("confirmPassword")
	c.Validation.Required(password != oldPassword).Message("输入了与旧密码相同的新密码").Key("password")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	user := c.forceCurrUser()
	c.Validation.Required(c.userService().VerifyPassword(user.CryptedPassword, oldPassword)).Message("你的旧密码输入有误").Key("oldPassword")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	if err := c.userService().DoChangePassword(&user, password); err != nil {
		c.Flash.Error("修改密码失败：" + err.Error())
	} else {
		c.Flash.Error("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
	}

	return c.Redirect(User.ChangePassword)
}

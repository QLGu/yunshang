package controllers

import (
	"strings"

	"github.com/revel/revel"
)

// 用户相关Actions
type User struct {
	ShouldLoginedController
}

// 用户主页
func (c User) Index() revel.Result {
	currUser, _ := c.currUser()

	userLevel, _ := c.userService().GetUserLevel(&currUser)
	revel.INFO.Printf("%v", userLevel)
	return c.Render(currUser, userLevel)
}

// 到用户信息
func (c User) UserInfo() revel.Result {
	return c.Render()
}

// 到修改密码
func (c User) ChangePassword() revel.Result {
	hasPassword := len(c.forceCurrUser().CryptedPassword) != 0
	return c.Render(hasPassword)
}

// 修改密码处理
func (c User) DoChangePassword(oldPassword, password, confirmPassword string) revel.Result {
	c.Validation.Required(oldPassword).Message("请输入旧密码")
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

// 到设置密码
func (c User) SetPassword(password, confirmPassword string) revel.Result {
	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	c.Validation.MaxSize(password, 50).Message("输入密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的密码不匹配").Key("confirmPassword")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	user := c.forceCurrUser()
	if err := c.userService().DoChangePassword(&user, password); err != nil {
		c.Flash.Error("修改密码失败：" + err.Error())
	} else {
		c.Flash.Error("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
	}

	return c.Redirect(User.ChangePassword)
}

// 用户级别显示
func (c User) ShowUserLevels() revel.Result {
	currUser, _ := c.currUser()
	userLevel, _ := c.userService().GetUserLevel(&currUser)

	userLevels := c.userService().FindUserLevels()
	userScores := currUser.Scores
	return c.Render(userLevels, userLevel, userScores)
}

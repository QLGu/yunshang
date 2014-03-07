package controllers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
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

	c.setChannel("user/index")
	return c.Render(currUser, userLevel)
}

// 到用户信息
func (c User) UserInfo() revel.Result {
	c.setChannel("user/userinfo")

	c.FlashParams()

	user := c.forceCurrUser()
	userDetail, _ := c.userService().GetUserDetailByUserId(user.Id)
	return c.Render(user, userDetail)
}

// 保存用户信息
func (c User) DoSaveUserInfo(user entity.User, userDetail entity.UserDetail) revel.Result {
	revel.INFO.Printf("user: %v", user)
	revel.INFO.Printf("userDetail: %v", userDetail)

	currUser := c.forceCurrUser()

	if len(currUser.LoginName) == 0 {
		c.Validation.Required(user.LoginName).Message("请输入登录名")
		c.Validation.MinSize(user.LoginName, 4).Message("请输入至少4位登录名")
		c.Validation.MaxSize(user.LoginName, 100).Message("输入过多位数的登录名")

		ok := true
		for _, v := range []rune(user.LoginName) {
			if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') || v == '_') {
				ok = false
				break
			}
		}
		c.Validation.Required(ok).Message("输入的登录名不符合要求")

		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}

		_, exists := c.userService().CheckUserByLoginName(user.LoginName)
		c.Validation.Required(!exists).Message("该用户名已经注册").Key("user.LoginName")
	}

	if len(currUser.Email) == 0 {
		revel.INFO.Println("len(currUser.Email) == 0 " + user.Email)
		c.Validation.Required(user.Email).Message("请输入邮箱")
		c.Validation.Email(user.Email).Message("请输入合法邮箱")
		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}
		_, exists := c.userService().CheckUserByEmail(user.Email)
		c.Validation.Required(!exists).Message("该邮箱已经注册").Key("user.Email")
		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}
	}

	err := c.userService().UpdateUserInfo(&currUser, user, userDetail)
	if err != nil {
		c.Flash.Error("保存用户信息失败" + err.Error())
	} else {
		c.Flash.Success("保存会员信息成功！")
	}
	return c.Redirect(User.UserInfo)
}

// 到修改密码
func (c User) ChangePassword() revel.Result {
	hasPassword := len(c.forceCurrUser().CryptedPassword) != 0

	c.setChannel("user/userinfo/cw")

	return c.Render(hasPassword)
}

// 修改密码处理
func (c User) DoChangePassword(oldPassword, password, confirmPassword string) revel.Result {
	c.Validation.Required(oldPassword).Message("请输入原始密码")
	c.Validation.Required(password).Message("请输入新密码")
	c.Validation.MinSize(password, 6).Message("请输入6位新密码")
	c.Validation.MaxSize(password, 50).Message("输入新密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认新密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的新密码不匹配").Key("confirmPassword")
	c.Validation.Required(password != oldPassword).Message("输入了与原始密码相同的新密码").Key("password")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	user := c.forceCurrUser()
	c.Validation.Required(c.userService().VerifyPassword(user.CryptedPassword, oldPassword)).Message("你的原始密码输入有误").Key("oldPassword")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	if err := c.userService().DoChangePassword(&user, password); err != nil {
		c.Flash.Error("修改密码失败：" + err.Error())
	} else {
		c.Flash.Success("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
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
		c.Flash.Success("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
	}

	return c.Redirect(User.ChangePassword)
}

// 用户级别显示
func (c User) UserLevel() revel.Result {
	currUser, _ := c.currUser()
	userLevel, _ := c.userService().GetUserLevel(&currUser)

	userLevels := c.userService().FindUserLevels()
	userScores := currUser.Scores

	c.setChannel("points/level")
	return c.Render(userLevels, userLevel, userScores)
}

// 积分规则显示
func (c User) ScoresRules() revel.Result {
	currUser, _ := c.currUser()
	userLevel, _ := c.userService().GetUserLevel(&currUser)

	userLevels := c.userService().FindUserLevels()
	userScores := currUser.Scores

	c.setChannel("points/rules")
	return c.Render(userLevels, userLevel, userScores)
}

// 用户的订单
func (c User) Orders() revel.Result {
	c.setChannel("order/orders")
	return c.Render()
}

// 显示用户头像
// param file： 头像图片标识： {{ucode}}.jpg
func (c User) Image(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.iamges", "data/images")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No Image Found！")
	}

	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

// 上传用户头像
func (c User) UploadImage(image *os.File) revel.Result {
	c.Validation.Required(image != nil)
	if ret := c.doValidate(User.UploadImage); ret != nil {
		return c.RenderJson(c.errorResposne("请选择图片", nil))
	}

	ucode, ok := c.Session["ucode"]
	gotang.Assert(ok, "ucode")
	gotang.Assert(len(ucode) != 0, "ucode")

	srcImage, err := imaging.Open(image.Name())
	if ret := c.checkError(err, "打开上传图片报错！"); ret != nil {
		return ret
	}

	dir := revel.Config.StringDefault("dir.data.iamges", "data/images")
	imageFile := filepath.Join(dir, ucode+".jpg")

	tnImage := imaging.Thumbnail(srcImage, 460, 460, imaging.Lanczos)
	err = imaging.Save(tnImage, imageFile)

	if ret := c.checkError(err, "保存上传图片报错！"); ret != nil {
		return ret
	}

	return c.RenderJson(c.successResposne("上传成功", nil))
}

func (c User) checkError(err error, msg string) revel.Result {
	if err != nil {
		revel.WARN.Printf("上传头像操作失败，%s， msg：%s", msg, err.Error())
		return c.RenderJson(c.errorResposne("操作失败，"+msg+", "+err.Error(), nil))
	}
	return nil
}

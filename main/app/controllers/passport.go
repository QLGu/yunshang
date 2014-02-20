package controllers

import (
	"github.com/dchest/captcha"
	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/robfig/revel"
)

type Passport struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
	AppController
}

func (c Passport) Reg(userType string) revel.Result {
	revel.INFO.Printf("userType: %v", userType)

	Captcha := struct {
		CaptchaId string
	}{
		captcha.New(),
	}

	return c.Render(Captcha)
}

func (c Passport) DoReg(userType, email, password, confirmPassword, validateCode, captchaId string) revel.Result {
	//c.Validation.Required(validateCode).Message("请输入验证码")
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误").Key("validateCode")

	c.Validation.Required(email).Message("请输入邮箱")
	c.Validation.Email(email).Message("请输入邮箱")
	c.Validation.MaxSize(email, 100).Message("输入邮箱位数太长了")

	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	c.Validation.MaxSize(password, 50).Message("输入密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的密码不匹配").Key("confirmPassword")

	if ret := c.checkReg(); ret != nil {
		return ret
	}
	c.Validation.Required(!c.userService().ExistsUserByEmail(email)).Message("%s邮箱已经被注册了", email).Key("email")
	if ret := c.checkReg(); ret != nil {
		return ret
	}

	user, err := c.userService().RegistUser(email, password)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(Passport.Reg)
	}

	data := struct {
		ActivationCode string
		Email          string
	}{user.ActivationCode, email}
	SendHtmlMail("激活邮件", utils.RenderTemplateToString("Passport/ActivateUserTemplate.html", data), email)

	c.RenderArgs["emailProvider"] = EmailProvider(email)
	c.RenderArgs["email"] = email
	return c.Render()
}

func (c Passport) Login(login string) revel.Result {
	Captcha := struct {
		CaptchaId string
		Login     string
	}{
		captcha.New(), login,
	}

	return c.Render(Captcha)
}

func (c Passport) DoLogin(login, password, validateCode, captchaId string) revel.Result {
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误").Key("validateCode")
	c.Validation.Required(login).Message("请输入账号")
	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	if ret := c.checkLogin(); ret != nil {
		return ret
	}

	user, ok := c.userService().CheckUser(login, password)
	c.Validation.Required(ok).Message("用户不存在或密码错误").Key("email")
	if ret := c.checkLogin(); ret != nil {
		return ret
	}

	// 执行登录后操作
	go c.userService().DoUserLogin(&user)

	c.Session["user"] = models.ToSessionUser(user).DisplayName()
	return c.Redirect(App.Index)
}

func (c Passport) Logout() revel.Result {
	delete(c.Session, "user")
	return c.Redirect(App.Index)
}

func (c Passport) Activate(activationCode string, email string) revel.Result {
	c.Validation.Email(email).Message("邮箱不合法!")

	revel.INFO.Printf("Activate Email: %v", email)
	if c.Validation.HasErrors() {
		c.RenderArgs["result"] = "邮箱不合法!"
		c.RenderArgs["activated"] = false
		return c.Render()
	}

	revel.INFO.Printf("activationCode: %v", activationCode)
	user, err := c.userService().Activate(email, activationCode)
	revel.INFO.Printf("Activate user: %v, enabled: %v", user, user.Enabled)

	opRet := ""
	if err != nil {
		opRet = err.Error()
	} else {
		opRet = "激活成功！"
		c.RenderArgs["activated"] = true
		c.RenderArgs["loginname"] = user.LoginName
		c.RenderArgs["email"] = user.Email
	}
	c.RenderArgs["result"] = opRet
	return c.Render()
}

func (c Passport) ForgotPasswordApply() revel.Result {
	Captcha := struct {
		CaptchaId string
	}{
		captcha.New(),
	}

	return c.Render(Captcha)
}

func (c Passport) DoForgotPasswordApply(email, validateCode, captchaId string) revel.Result {
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误").Key("validateCode")
	c.Validation.Email(email).Message("请输入合法的邮箱")
	if ret := c.check(Passport.ForgotPasswordApply); ret != nil {
		return ret
	}

	user, ok := c.userService().CheckUserByEmail(email)
	c.Validation.Required(ok).Message("请输入你注册的邮箱").Key("email")
	if ret := c.check(Passport.ForgotPasswordApply); ret != nil {
		return ret
	}

	c.userService().DoForgotPasswordApply(&user)

	data := struct {
		PasswordResetCode string
		Email             string
	}{user.PasswordResetCode, email}
	SendHtmlMail("激活邮件", utils.RenderTemplateToString("Passport/ResetPasswordTemplate.html", data), user.Email)

	c.RenderArgs["emailProvider"] = EmailProvider(email)
	c.RenderArgs["email"] = email

	return c.Render()
}

func (c Passport) DoResetPassword(email, passwordResetCode string) revel.Result {
	c.Validation.Email(email).Message("请输入合法的邮箱")
	if c.Validation.HasErrors() {
		c.RenderArgs["result"] = "输入不合法"
		c.RenderArgs["ok"] = false
		return c.Render()
	}

	newPassword, err := c.userService().ResetUserPassword(email, passwordResetCode)
	if err != nil {
		c.RenderArgs["result"] = err.Error()
		c.RenderArgs["ok"] = false
		return c.Render()
	}

	c.RenderArgs["email"] = email
	c.RenderArgs["ok"] = true
	return c.Render(newPassword)
}

///////////////////////////////////////////////////////////////////////////////////

func (c Passport) checkReg() revel.Result {
	return c.check(Passport.Reg)
}

func (c Passport) checkLogin() revel.Result {
	return c.check(Passport.Login)
}

func (c Passport) check(f interface{}) revel.Result {
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(f)
	}
	return nil
}

func (c Passport) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}
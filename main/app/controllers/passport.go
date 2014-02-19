package controllers

import (
	"fmt"

	"github.com/dchest/captcha"
	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/models"
	"github.com/robfig/revel"
)

type Passport struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
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
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误")

	c.Validation.Required(email).Message("请输入邮箱")
	c.Validation.Email(email).Message("请输入邮箱")

	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	c.Validation.Required(confirmPassword).Message("请输入确认密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的密码不匹配")

	if ret := c.checkReg(); ret != nil {
		return ret
	}
	c.Validation.Required(!c.userService().ExistsUserByEmail(email)).Message("%s邮箱已经被注册了", email)
	if ret := c.checkReg(); ret != nil {
		return ret
	}

	user, err := c.userService().RegistUser(email, password)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(Passport.Reg)
	}

	// TODO 邮件内容模板
	SendHtmlMail("激活邮件", fmt.Sprintf(`
	<a href="http://localhost:9000/passport/activate?activationCode=%s&email=%s">激活</a>
	`, user.ActivationCode, email), email)

	c.RenderArgs["emailProvider"] = EmailProvider(email)
	c.RenderArgs["email"] = email
	return c.Render()
}

func (c Passport) Login() revel.Result {
	return c.Render()
}

func (c Passport) DoLogin(email string) revel.Result {
	revel.INFO.Printf("email:%v", email)

	c.Session["user"] = email
	return c.Render()
}

func (c Passport) Logout() revel.Result {
	delete(c.Session, "user")
	return c.Redirect(App.Index)
}

func (c Passport) Activate(activationCode string, email string) revel.Result {
	revel.INFO.Printf("activationCode: %v", activationCode)
	user, err := c.userService().Activate(email, activationCode)
	revel.INFO.Printf("Activate user: %v, enabled: %v", user, user.Enabled)
	if err != nil {
		c.Flash.Error(err.Error())
	} else {
		c.Flash.Success("激活成功！")
		c.RenderArgs["activated"] = true
		c.RenderArgs["loginname"] = user.LoginName
		c.RenderArgs["email"] = user.Email
	}
	return c.Render()
}

func (c Passport) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c Passport) checkReg() revel.Result {
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Passport.Reg)
	}
	return nil
}

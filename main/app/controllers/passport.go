package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/dchest/captcha"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/mail"
	"github.com/itang/yunshang/modules/oauth"
	"github.com/revel/revel"
)

// 登录相关的Actions
type Passport struct {
	AppController
}

// 注册
func (c Passport) Reg(userType string) revel.Result {
	captchaId := captcha.NewLen(4)

	c.setChannel("index/reg")
	return c.Render(captchaId)
}

// 注册处理
func (c Passport) DoReg(userType, email, password, confirmPassword, validateCode, captchaId string) revel.Result {
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
	c.Validation.Required(!c.userApi().ExistsUserByEmail(email)).Message("%s邮箱已经被注册了", email).Key("email")
	if ret := c.checkReg(); ret != nil {
		return ret
	}

	user, err := c.userApi().RegistUser(email, password)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(Passport.Reg)
	}

	err = gotang.DoIOWithTimeout(func() error {
		return mail.SendHtml("激活邮件",
			utils.RenderTemplate("Passport/ActivateUserTemplate.html", struct {
				ActivationCode string
				Email          string
			}{user.ActivationCode, email}),
			email)
	}, time.Second*60)
	if err != nil { // TODO
		revel.ERROR.Println("发送邮件超时", err)

		c.setRollbackOnly()
		c.Flash.Error("注册处理超时了，请重试！")
		return c.Redirect(Passport.Reg)
	}

	c.RenderArgs["emailProvider"] = mail.GetEmailProvider(email)
	c.RenderArgs["email"] = email
	return c.Render()
}

// 到登录
func (c Passport) Login(login string) revel.Result {
	captchaId := captcha.NewLen(4)
	providers := oauth.GetProviders()

	c.setChannel("index/login")
	return c.Render(login, captchaId, providers)
}

// 登录处理
func (c Passport) DoLogin(login, password, validateCode, captchaId string) revel.Result {
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误").Key("validateCode")
	c.Validation.Required(login).Message("请输入账号")
	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	if ret := c.checkLogin(); ret != nil {
		return ret
	}

	user, ok := c.userApi().CheckUser(login, password)
	c.Validation.Required(ok).Message("用户不存在或密码错误或未激活。有任何疑问，请联系本站客户！").Key("login")
	if ret := c.checkLogin(); ret != nil {
		return ret
	}

	// 执行登录后操作
	go c.userApi().DoUserLogin(&user)

	c.setLoginSession(models.ToSessionUser(user))

	return c.Redirect(App.Index)
}

func (c Passport) DoLoginFromIndex(login, password string) revel.Result {
	c.Validation.Required(login).Message("请输入账号")
	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	if ret := c.checkLogin(); ret != nil {
		return c.RenderJson(Error("输入有误!", nil))
	}

	user, ok := c.userApi().CheckUser(login, password)
	c.Validation.Required(ok).Message("用户不存在或密码错误或未激活。有任何疑问，请联系本站客户！").Key("login")
	if ret := c.checkLogin(); ret != nil {
		return c.RenderJson(Error("输入有误!", nil))
	}

	// 执行登录后操作
	go c.userApi().DoUserLogin(&user)

	c.setLoginSession(models.ToSessionUser(user))

	return c.RenderJson(Success("登录成功!", nil))
}

// 开放平台登录入口
func (c Passport) OpenLogin(provider string) revel.Result {
	return SocialAuth.HandleRedirect(c.Controller)
}

// 开放平台登录处理
func (c Passport) DoOpenLogin(code string) revel.Result {
	log.Printf("%v", c.Params)
	log.Printf("weibo code: " + code)

	redirect, _, err := SocialAuth.OAuthAccess(c.Controller, c.db)
	if err != nil {
		revel.ERROR.Printf("SocialAuth.handleAccess, %v", err)
	}

	return c.Redirect(redirect)
}

// 连接第三方用户
func (c Passport) Connect() revel.Result {
	st, ok := SocialAuth.ReadyConnect(c.Controller)
	if !ok {
		return c.Redirect(Passport.Login)
	}

	// Your app need custom connect behavior
	// example just direct connect and login
	//TODO: 关联用户ID??
	p, stf, err := c.getProviderAndTokenFromSession()
	if err != nil {
		revel.ERROR.Printf("getProviderAndTokenFromSession, %v", err)
		c.Flash.Error(fmt.Sprintf("getProviderAndTokenFromSession: %v", err))
		return c.Redirect(Passport.Login)
	}

	identify, err := p.GetIndentify(stf.Token)
	if err != nil {
		revel.ERROR.Printf("getProviderAndTokenFromSession, %v", err)
		c.Flash.Error(fmt.Sprintf("GetIndentify: %v", err))
		return c.Redirect(Passport.Login)
	}

	email := "" // TODO get from provider
	user, err := c.userApi().ConnectUser(identify, st.Name(), email)
	if err != nil {
		revel.ERROR.Printf("ConnectUser, %v", err)
		c.Flash.Error(fmt.Sprintf("ConnectUser: %v", err))
		return c.Redirect(Passport.Login)
	}

	uid := user.Id
	loginRedirect, _, err := SocialAuth.ConnectAndLogin(c.db, c.Controller, st, uid)
	if err != nil {
		revel.ERROR.Printf("SocialAuth.handleAccess, %v", err)
	}

	return c.Redirect(loginRedirect)
}

// 退出登录
func (c Passport) Logout() revel.Result {
	c.clearLoginSession()

	return c.Redirect(App.Index)
}

// 到激活
func (c Passport) Activate(activationCode string, email string) revel.Result {
	c.Validation.Email(email).Message("邮箱不合法!")

	revel.INFO.Printf("Activate Email: %v", email)
	if c.Validation.HasErrors() {
		c.RenderArgs["result"] = "邮箱不合法!"
		c.RenderArgs["activated"] = false
		return c.Render()
	}

	revel.INFO.Printf("activationCode: %v", activationCode)
	user, err := c.userApi().Activate(email, activationCode)
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

// 到忘记密码申请
func (c Passport) ForgotPasswordApply() revel.Result {
	captchaId := captcha.NewLen(4)

	c.setChannel("index/forget_pwd")
	return c.Render(captchaId)
}

// 忘记密码处理
func (c Passport) DoForgotPasswordApply(email, validateCode, captchaId string) revel.Result {
	c.Validation.Required(captcha.VerifyString(captchaId, validateCode)).Message("验证码填写有误").Key("validateCode")
	c.Validation.Email(email).Message("请输入合法的邮箱")
	if ret := c.doValidate(Passport.ForgotPasswordApply); ret != nil {
		return ret
	}

	user, ok := c.userApi().CheckUserByEmail(email)
	c.Validation.Required(ok).Message("请输入你注册的邮箱").Key("email")
	if ret := c.doValidate(Passport.ForgotPasswordApply); ret != nil {
		return ret
	}

	c.userApi().DoForgotPasswordApply(&user)

	err := gotang.DoIOWithTimeout(func() error {
		return mail.SendHtml("重置密码邮件", utils.RenderTemplate("Passport/ResetPasswordTemplate.html", struct {
			PasswordResetCode string
			Email             string
		}{user.PasswordResetCode, email}), user.Email)
	}, time.Second*60)
	if err != nil {
		revel.ERROR.Println("发送邮件超时", err)

		c.setRollbackOnly()
		c.Flash.Error("忘记密码申请处理超时了，请重试！")
		return c.Redirect(Passport.ForgotPasswordApply)
	}

	c.RenderArgs["emailProvider"] = mail.GetEmailProvider(email)
	c.RenderArgs["email"] = email

	return c.Render()
}

// 重置密码处理
func (c Passport) DoResetPassword(email, passwordResetCode string) revel.Result {
	c.Validation.Email(email).Message("请输入合法的邮箱")
	if c.Validation.HasErrors() {
		c.RenderArgs["result"] = "输入不合法"
		c.RenderArgs["ok"] = false
		return c.Render()
	}

	newPassword, err := c.userApi().ResetUserPassword(email, passwordResetCode)
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
// private

func (c Passport) checkReg() revel.Result {
	return c.doValidate(Passport.Reg)
}

func (c Passport) checkLogin() revel.Result {
	return c.doValidate(Passport.Login)
}

func (c Passport) getProviderAndTokenFromSession() (p oauth.Provider, stf oauth.SocialTokenField, err error) {
	socialType, ok := SocialAuth.ReadyConnect(c.Controller)
	if !ok {
		return p, stf, fmt.Errorf("can't ReadyConnect")
	}
	tokKey := SocialAuth.GetSessKey(socialType, "token")

	tk := oauth.SocialTokenField{}
	value := c.Session[tokKey]
	if err := tk.SetRaw(value); err != nil {
		return p, stf, err
	}

	if p, _ = oauth.GetProviderByType(socialType); p == nil {
		return p, stf, fmt.Errorf("unknown provider")
	}

	return p, tk, nil
}

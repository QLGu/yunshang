package controllers

import (
	"github.com/robfig/revel"

	reveltang "github.com/itang/reveltang/controllers"
)

type Passport struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
}

func (c Passport) Reg(userType string) revel.Result {
	revel.INFO.Printf("userType: %v", userType)

	return c.Render()
}

func (c Passport) DoReg(userType string) revel.Result {
	c.RenderArgs["emailProvider"] = "http://mail.google.com"
	c.RenderArgs["email"] = "live.tang@gmail.com"
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

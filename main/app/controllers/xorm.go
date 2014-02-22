package controllers

import (
	"github.com/lunny/xorm"
	"github.com/robfig/revel"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app"
)

type XOrmController struct {
	*revel.Controller
	Engine *xorm.Engine
}

type XOrmTnController struct {
	*revel.Controller
	XOrmSession *xorm.Session
}

func (self *XOrmController) begin() revel.Result {
	self.Engine = app.Engine
	return nil
}

func (self *XOrmTnController) begin() revel.Result {
	gotang.Assert(app.Engine != nil, "app.Engine can't be nil")

	self.XOrmSession = app.Engine.NewSession()
	self.XOrmSession.Begin()

	return nil
}

func (self *XOrmTnController) commit() revel.Result {
	self.XOrmSession.Commit()
	self.XOrmSession.Close()

	return nil
}

func (self *XOrmTnController) rollback() revel.Result {
	self.XOrmSession.Rollback()
	self.XOrmSession.Close()

	return nil
}

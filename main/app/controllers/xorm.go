package controllers

import (
	"github.com/itang/gotang"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// xorm事务控制
type XOrmController struct {
	*revel.Controller
	Engine *xorm.Engine
}

type XOrmTnController struct {
	*revel.Controller
	XOrmSession *xorm.Session
}

func (self *XOrmController) begin() revel.Result {
	gotang.Assert(db.Engine != nil, "db.Engine can't be nil")

	self.Engine = db.Engine
	return nil
}

func (self *XOrmTnController) begin() revel.Result {
	gotang.Assert(db.Engine != nil, "db.Engine can't be nil")

	self.XOrmSession = db.Engine.NewSession()
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

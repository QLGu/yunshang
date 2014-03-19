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
	db *xorm.Engine
}

type XOrmTnController struct {
	*revel.Controller
	db *xorm.Session
}

func (self *XOrmController) begin() revel.Result {
	gotang.Assert(db.Engine != nil, "db.Engine can't be nil")

	self.db = db.Engine
	self.RenderArgs["_db"] = self.db
	return nil
}

func (self *XOrmTnController) begin() revel.Result {
	gotang.Assert(db.Engine != nil, "db.Engine can't be nil")

	self.db = db.Engine.NewSession()
	self.db.Begin()

	self.RenderArgs["_db"] = self.db

	return nil
}

func (self *XOrmTnController) commit() revel.Result {
	self.db.Commit()
	self.db.Close()

	return nil
}

func (self *XOrmTnController) rollback() revel.Result {
	self.db.Rollback()
	self.db.Close()

	return nil
}

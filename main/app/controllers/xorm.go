package controllers

import (
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/utils"
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
	db           *xorm.Session
	rollbackOnly bool
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

	//events
	self.db.After(func(bean interface{}) {
		models.Emitter.Emit(models.EUpdateCache, utils.TypeOfTarget(bean).Name())
		//fmt.Println("after, ", bean, "name:", utils.TypeOfTarget(bean), utils.TypeOfTarget(bean).Name());
	})

	self.db.Begin()

	self.RenderArgs["_db"] = self.db

	return nil
}

func (self *XOrmTnController) commit() revel.Result {
	if self.rollbackOnly {
		return self.rollback()
	}

	self.db.Commit()
	self.db.Close()

	return nil
}

func (self *XOrmTnController) rollback() revel.Result {
	self.db.Rollback()
	self.db.Close()

	return nil
}

func (self *XOrmTnController) setRollbackOnly() {
	self.rollbackOnly = true
}

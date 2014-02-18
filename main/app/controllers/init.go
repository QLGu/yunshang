package controllers

import (
	"github.com/robfig/revel"
)

func init() {
	revel.InterceptMethod((*XOrmController).Begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).Begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).Commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).Rollback, revel.PANIC)
}

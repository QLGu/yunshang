package controllers

import (
	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

func init() {
	revel.InterceptMethod((*XOrmController).Begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).Begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).Commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).Rollback, revel.PANIC)

	email.InitGmail("live.tang@gmail.com", "tq19811115")
	email.Config.From.Name = "YuShang"
}

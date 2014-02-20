package controllers

import (
	"github.com/robfig/revel"
	"github.com/ungerik/go-mail"
)

func init() {
	revel.InterceptMethod((*XOrmController).begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).rollback, revel.PANIC)

	revel.InterceptMethod((*ShouldLoginedController).checkUser, revel.BEFORE)
	revel.InterceptMethod((*AdminController).checkAdminUser, revel.BEFORE)

	email.InitGmail("live.tang@gmail.com", "tq19811115")
	email.Config.From.Name = "YuShang"
}

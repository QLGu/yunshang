package controllers

import (
	//"fmt"
	//"time"

	//"github.com/itang/gotang"
	//"github.com/itang/yunshang/main/app/models/entity"
	//"github.com/itang/yunshang/main/app/utils"
	//"github.com/itang/yunshang/modules/mail"
	//"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 产品相关Actions
type Product struct {
	AppController
}

func (c Product) View(id int64) revel.Result {
	return c.Render()
}

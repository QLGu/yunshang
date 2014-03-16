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

// 制造商相关Actions
type Provider struct {
	AppController
}

func (c Provider) View(id int64) revel.Result {
	if id == 0 {
		return c.NotFound("制造商不存在！")
	}

	p, ok := c.productApi().GetProviderById(id)
	if !ok {
		return c.NotFound("制造商不存在！")
	}
	return c.Render(p)
}

func (c Provider) ProvidersData() revel.Result {
	ps := c.productApi().FindAllAvailableProviders()
	return c.RenderJson(c.successResposne("", ps))
}

func (c Provider) ProviderData(id int64) revel.Result {
	revel.INFO.Println("id", id)

	p, _ := c.productApi().GetProviderById(id)
	return c.RenderJson(c.successResposne("", p))
}

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

type P struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (c Product) ProvidersData() revel.Result {
	ps := []P{{1, "西门子电子"}, {2, "华立LED"}}
	return c.RenderJson(c.successResposne("dd", ps))
}

func (c Product) ProviderData(id int64) revel.Result {
	revel.INFO.Println("id", id)

	ps := []P{{1, "西门子电子"}, {2, "华立LED"}}
	return c.RenderJson(c.successResposne("dd", ps[id-1]))
}

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

// 商品分类相关Actions
type Category struct {
	AppController
}

func (c Category) CategoriesData() revel.Result {
	ps := c.productApi().FindAllAvailableCategories()
	return c.RenderJson(c.successResposne("", ps))
}

func (c Category) CategoryData(id int64) revel.Result {
	p, _ := c.productApi().GetCategoryById(id)
	return c.RenderJson(c.successResposne("", p))
}
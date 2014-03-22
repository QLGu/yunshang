package controllers

import (
	"github.com/revel/revel"
)

// 商品分类相关Actions
type Category struct {
	AppController
}

func (c Category) Index() revel.Result {
	c.setChannel("categories/index")
	cgs := c.productApi().FindAvailableTopCategories()
	return c.Render(cgs)
}

func (c Category) CategoriesData() revel.Result {
	ps := c.productApi().FindAllAvailableCategories()
	return c.RenderJson(Success("", ps))
}

func (c Category) CategoryData(id int64) revel.Result {
	p, _ := c.productApi().GetCategoryById(id)
	return c.RenderJson(Success("", p))
}

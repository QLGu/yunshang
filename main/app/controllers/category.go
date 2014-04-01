package controllers

import (
	"github.com/revel/revel"
)

// 商品分类相关Actions
type Category struct {
	AppController
}

func (c Category) Index() revel.Result {
	c.setChannel("products/categories")
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

func (c Category) NewsCategoriesData() revel.Result {
	ps := c.newsApi().FindAllAvailableCategories()
	return c.RenderJson(Success("", ps))
}

func (c Category) NewsCategoryData(id int64) revel.Result {
	p, _ := c.newsApi().GetCategoryById(id)
	return c.RenderJson(Success("", p))
}

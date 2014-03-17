package controllers

import (
	"github.com/revel/revel"
)

// 相关Actions
type About struct {
	AppController
}

// 主页
func (c About) Index() revel.Result {
	c.setChannel("about/index")

	products := c.productApi().FindAllAvailableProducts()
	categories := c.productApi().FindAvailableTopCategories()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, categories, providers)
}

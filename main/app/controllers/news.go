package controllers

import (
	"github.com/revel/revel"
)

// 相关Actions
type News struct {
	AppController
}

// 产品主页
func (c News) Index() revel.Result {
	c.setChannel("news/index")

	products := c.productApi().FindAllAvailableProducts()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, providers)
}

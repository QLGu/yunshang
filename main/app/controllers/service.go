package controllers

import (
	"github.com/revel/revel"
)

// 产品相关Actions
type Service struct {
	AppController
}

// 产品主页
func (c Service) Index() revel.Result {
	c.setChannel("services/index")

	products := c.productApi().FindAllAvailableProducts()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, providers)
}

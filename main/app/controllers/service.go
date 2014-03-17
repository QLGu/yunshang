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
	categories := c.productApi().FindAvailableTopCategories()
	return c.Render(categories)
}

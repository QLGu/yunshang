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
	categories := c.productApi().FindAvailableTopCategories()
	return c.Render(categories)
}

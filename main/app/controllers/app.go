package controllers

import (
	"github.com/revel/revel"
)

// 应用主控制器
type App struct {
	AppController
}

// 应用主页
func (c App) Index() revel.Result {
	c.setChannel("index/")
	products := c.productApi().FindAllAvailableProducts()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, providers)
}

func (c App) AdImagesData() revel.Result {
	images := c.appApi().FindAdImages()
	return c.RenderJson(Success("", images))

}

func (c App) AdImage(file string) revel.Result {
	targetFile, err := c.appApi().GetAdImageFile(file)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

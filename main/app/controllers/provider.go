package controllers

import (
	"os"
	"path/filepath"

	gio "github.com/itang/gotang/io"

	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 制造商相关Actions
type Provider struct {
	AppController
}

// 品牌主页
func (c Provider) Index() revel.Result {
	c.setChannel("providers/index")

	products := c.productApi().FindAllAvailableProducts()
	categories := c.productApi().FindAvailableTopCategories()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, categories, providers)
}

func (c Provider) View(id int64) revel.Result {
	if id == 0 {
		return c.NotFound("制造商不存在！")
	}

	p, exists := c.productApi().GetProviderById(id)
	if !exists {
		return c.NotFound("制造商不存在！")
	}
	c.setChannel("providers/view")

	products := c.productApi().FindAllAvailableProducts()
	categories := c.productApi().FindAvailableTopCategories()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(p, products, categories, providers)
}

func (c Provider) ProvidersData() revel.Result {
	ps := c.productApi().FindAllAvailableProviders()
	return c.RenderJson(Success("", ps))
}

func (c Provider) ProviderData(id int64) revel.Result {
	revel.INFO.Println("id", id)

	p, _ := c.productApi().GetProviderById(id)
	return c.RenderJson(Success("", p))
}

// param file： 头像图片标识： {{id}}.jpg
func (c Provider) Image(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.providers", "data/providers")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "p-default.jpg")
	}

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No Image Found！")
	}

	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

// 产品数据
func (c Provider) ProductData(providerId int64) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		session.And("provider_id=?", providerId)
	})
	page := c.productApi().FindAllAvailableProductsForPage(ps)
	return c.renderDTJson(page)
}

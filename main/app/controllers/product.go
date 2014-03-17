package controllers

import (
	"fmt"
	"os"
	"path/filepath"

	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/revel/revel"
)

// 产品相关Actions
type Product struct {
	AppController
}

// 产品主页
func (c Product) Index() revel.Result {
	c.setChannel("products/index")

	products := c.productApi().FindAllAvailableProducts()
	categories := c.productApi().FindAvailableTopCategories()
	providers := c.productApi().FindAllAvailableProviders()
	return c.Render(products, categories, providers)
}

func (c Product) View(id int64) revel.Result {
	if id == 0 {
		return c.NotFound("产品不存在！")
	}
	p, exists := c.productApi().GetProductById(id)
	if !exists {
		return c.NotFound("产品不存在！")
	}

	provider, _ := c.productApi().GetProviderByProductId(p.Id)
	detail, _ := c.productApi().GetProductDetail(p.Id)

	products := c.productApi().FindAllAvailableProducts()
	categories := c.productApi().FindAvailableTopCategories()
	providers := c.productApi().FindAllAvailableProviders()

	c.setChannel("products/view")
	return c.Render(p, provider, detail, products, categories, providers)
}

func (c Product) SdImages(id int64) revel.Result {
	images := c.productApi().FindProductImages(id, entity.PTScheDiag)
	return c.RenderJson(c.successResposne("", images))
}

func (c Product) SdImage(file string) revel.Result {
	targetFile, err := c.productApi().GetProductImageFile(file, entity.PTScheDiag)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

func (c Product) ImagePics(id int64) revel.Result {
	images := c.productApi().FindProductImages(id, entity.PTPics)
	return c.RenderJson(c.successResposne("", images))
}

func (c Product) ImagePicsList(id int64) revel.Result {
	var images []entity.ProductParams
	c.XOrmSession.Where("type=? and product_id=?", entity.PTPics, id).Find(&images)
	var ret = ""
	for _, v := range images {
		ret += fmt.Sprintf("?file=%sue_separate_ue", v.Value)
	}
	return c.RenderText(ret)
}

func (c Product) ImagePic(file string) revel.Result {
	targetFile, err := c.productApi().GetProductImageFile(file, entity.PTPics)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

func (c Product) MFiles(id int64) revel.Result {
	var files []entity.ProductParams
	c.XOrmSession.Where("type=? and product_id=?", entity.PTMaterial, id).Find(&files)
	return c.RenderJson(c.successResposne("", files))
}

// 材料
func (c Product) MFile(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.product.m", "data/products/m")
	f := filepath.Join(dir, filepath.Base(file))
	targetFile, err := os.Open(f)
	if err != nil {
		return c.NotFound("No File Found！")
	}

	return c.RenderFile(targetFile, "")
}

func (c Product) Specs(id int64) revel.Result {
	var files []entity.ProductParams
	c.XOrmSession.Where("type=? and product_id=?", entity.PTSpec, id).Find(&files)
	return c.RenderJson(c.successResposne("", files))
}

func (c Product) Prices(id int64) revel.Result {
	var prices []entity.ProductPrices
	c.XOrmSession.Where("product_id=?", id).Find(&prices)
	return c.RenderJson(c.successResposne("", prices))
}

// param file： 头像图片标识： {{id}}.jpg
func (c Product) Image(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.products.logo", "data/products/logo")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No Image Found！")
	}

	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

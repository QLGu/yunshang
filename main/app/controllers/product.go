package controllers

import (
	//"fmt"
	//"time"
	"os"
	"path/filepath"

	//"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	//"github.com/itang/yunshang/main/app/utils"
	//"github.com/itang/yunshang/modules/mail"
	gio "github.com/itang/gotang/io"
	"github.com/revel/revel"
)

// 产品相关Actions
type Product struct {
	AppController
}

func (c Product) View(id int64) revel.Result {
	p, _ := c.productApi().GetProductById(id)
	provider, _ := c.productApi().GetProviderByProductId(p.Id)
	return c.Render(p, provider)
}

func (c Product) SdImages(id int64) revel.Result {
	var images []entity.ProductParams
	c.XOrmSession.Where("type=? and product_id=?", entity.PTScheDiag, id).Find(&images)
	return c.RenderJson(c.successResposne("", images))
}

func (c Product) SdImage(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.product.sd", "data/products/sd")

	imageFile := filepath.Join(dir, filepath.Base(file))
	revel.INFO.Println(imageFile)
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

func (c Product) MFiles(id int64) revel.Result {
	var files []entity.ProductParams
	c.XOrmSession.Where("type=? and product_id=?", entity.PTMaterial, id).Find(&files)
	return c.RenderJson(c.successResposne("", files))
}

// 材料
func (c Product) MFile(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.product.m", "data/products/m")

	imageFile := filepath.Join(dir, filepath.Base(file))

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No File Found！")
	}

	return c.RenderFile(targetFile, "")
}

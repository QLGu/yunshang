package controllers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/revel/revel"
)

// 产品相关Actions
type Product struct {
	AppController
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
	return c.Render(p, provider, detail)
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

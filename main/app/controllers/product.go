package controllers

import (
	"fmt"
	"os"
	"path/filepath"

	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 产品相关Actions
type Product struct {
	AppController
}

// 产品主页
func (c Product) Index(ctcode string, p int64, q string, hide_filters string) revel.Result {
	c.setChannel("products/index")

	//当前查询对应的制造商
	var providers = make([]entity.Provider, 0)
	if p != 0 {
		p, exists := c.productApi().GetProviderById(p)
		if exists {
			providers = append(providers, p)
		}
	}

	//当前查询对应的分类
	pcts := c.productApi().FindAvailableCategoryChainByCode(ctcode)

	// 分类过滤条件
	filterCts := c.productApi().RecommendCategories()
	//制造商过滤条件
	filterPs := c.productApi().RecommendProviders()

	// 查询产品结果数据
	ps := c.pageSearcherWithCalls(func(s *xorm.Session) {
		if len(ctcode) > 0 {
			s.And("category_code ilike ?", ctcode+"%")
		}
		if p != 0 {
			s.And("provider_id=?", p)
		}

		//名称和型号
		if len(q) > 0 {
			s.And("(name ilike ? or model ilike ?)", "%"+q+"%", "%"+q+"%")
		}
	})

	//指定数量: 5 * 6
	if(ps.Limit ==10){
		ps.Limit = 30
		c.RenderArgs["limit"] = ps.Limit
	}
	products := c.productApi().FindAllAvailableProductsForPage(ps)

	return c.Render(ctcode, p, q, hide_filters, pcts, providers, filterPs, filterCts, products)
}

func (c Product) ProvidersForSelect() revel.Result {
	ps := c.productApi().FindProvidersForSearchSelect()
	return c.Render(ps)
}

func (c Product) CategoriesForSelect() revel.Result {
	cs := c.productApi().FindCategoriesForSearchSelect()
	return c.Render(cs)
}

func (c Product) IndexByCategory(code string) revel.Result {
	return c.Redirect(routes.Product.Index(code, 0, "", ""))
}

func (c Product) View(id int64) revel.Result {
	if id == 0 {
		return c.NotFound("产品不存在！")
	}
	p, exists := c.productApi().GetProductById(id)
	if !exists {
		return c.NotFound("产品不存在！")
	}

	//供应商
	provider, _ := c.productApi().GetProviderByProductId(p.Id)
	//产品详情
	detail, _ := c.productApi().GetProductDetail(p.Id)

	c.setChannel("products/view")
	return c.Render(p, provider, detail)
}

//产品评论
func (c Product) ProductComments(id int64) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		session.And("enabled=?", true)
		session.And("target_type=?", entity.CT_PRODUCT)
		session.And("target_id=?", id)
	})
	pageObject := c.userApi().CommentsForPage(ps)
	return c.Render(id, pageObject)
}

func (c Product) SdImages(id int64) revel.Result {
	images := c.productApi().FindProductImages(id, entity.PTScheDiag)
	return c.RenderJson(Success("", images))
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
	return c.RenderJson(Success("", images))
}

func (c Product) ImagePicsList(id int64) revel.Result {
	var images []entity.ProductParams
	c.db.Where("type=? and product_id=?", entity.PTPics, id).Find(&images)
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
	c.db.Where("type=? and product_id=?", entity.PTMaterial, id).Find(&files)
	return c.RenderJson(Success("", files))
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
	c.db.Where("type=? and product_id=?", entity.PTSpec, id).Find(&files)
	return c.RenderJson(Success("", files))
}

func (c Product) Prices(id int64) revel.Result {
	prices := c.productApi().FindProductPrices(id)
	return c.RenderJson(Success("", prices))
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

package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

type ProductService interface {
	Total() int64

	FindAllAvailableProducts() []entity.Product

	FindAllProductsForPage(ps *PageSearcher) PageData

	SaveProduct(p entity.Product) (id int64, err error)

	GetProductById(id int64) (entity.Product, bool)

	ToggleProductEnabled(p *entity.Product) error

	SaveProductDetail(id int64, content string) error

	GetProductDetail(id int64) (detail string, err error)

	GetProductImageFile(file string, t int) (*os.File, error)

	FindProductImages(id int64, t int) []entity.ProductParams

	DeleteProductParams(id int64) error

	FindAllProvidersForPage(ps *PageSearcher) PageData

	FindAllAvailableProviders() []entity.Provider

	SaveProvider(p entity.Provider) (id int64, err error)

	GetProviderById(id int64) (entity.Provider, bool)

	GetProviderByProductId(id int64) (entity.Provider, bool)

	ToggleProviderEnabled(p *entity.Provider) error

	DeleteProvider(id int64) error

	FindAllCategoriesForPage(ps *PageSearcher) PageData

	FindAllAvailableCategories() []entity.ProductCategory

	FindAvailableTopCategories() []entity.ProductCategory

	FindAllAvailableCategoriesByParentId(id int64) []entity.ProductCategory

	SaveCategory(p entity.ProductCategory) (id int64, err error)

	GetCategoryById(id int64) (entity.ProductCategory, bool)

	ToggleCategoryEnabled(p *entity.ProductCategory) error

	FindAllProductStockLogs(id int64) []entity.ProductStockLog

	AddProductStock(id int64, stock int, message string) (int, error)
}

func NewProductService(session *xorm.Session) ProductService {
	return &productService{session}
}

/////////////////////////////////////////////////////////////////////////////

type productService struct {
	session *xorm.Session
}

func (self productService) Total() int64 {
	total, _ := self.session.Count(&entity.Product{})
	return total
}

func (self productService) FindAllProductsForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(session *xorm.Session) {
		session.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Product{})
	if err != nil {
		log.Println(err)
	}

	var products []entity.Product

	err1 := ps.BuildQuerySession().Find(&products)
	if err1 != nil {
		log.Println(err1)
	}

	return NewPageData(total, products)
}

func (self productService) FindAllAvailableProducts() (ps []entity.Product) {
	_ = self.session.Where("enabled=?", true).Find(&ps)
	return
}

func (self productService) GetProductById(id int64) (p entity.Product, ok bool) {
	ok, _ = self.session.Where("id=?", id).Get(&p)
	return
}

func (self productService) SaveProduct(p entity.Product) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.session.Insert(&p)
		if err != nil {
			return
		}
		p.Code = p.Id + entity.ProductStartDisplayCode //编码
		_, err = self.session.Id(p.Id).Cols("code").Update(&p)
		id = p.Id

		Emitter.Emit(EStockLog, "", p.Id, fmt.Sprintf("创建产品入库：%d(%s)", p.StockNumber, p.UnitName))

		return
	} else { // update
		currDa, ok := self.GetProductById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.session.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此产品不存在")
		}
	}
}

func (self productService) AddProductStock(id int64, stock int, message string) (newstock int, err error) {
	p, ok := self.GetProductById(id)
	if !ok {
		err = errors.New("产品不存在!")
		return
	}
	p.StockNumber += stock
	_, err = self.session.Id(p.Id).Cols("stock_number").Update(&p)
	Emitter.Emit(EStockLog, "", p.Id, fmt.Sprintf("入库：%d(%s), 当前库存%d(%s)", stock, p.UnitName, p.StockNumber, p.UnitName))

	newstock = p.StockNumber
	return
}

func (self productService) ToggleProductEnabled(p *entity.Product) error {
	p.Enabled = !p.Enabled
	if p.Enabled {
		p.EnabledAt = time.Now()
		_, err := self.session.Id(p.Id).Cols("enabled", "enabled_at").Update(p)
		return err
	}
	p.UnEnabledAt = time.Now()
	_, err := self.session.Id(p.Id).Cols("enabled", "un_enabled_at").Update(p)
	return err
}

func (self productService) FindAllProvidersForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(session *xorm.Session) {
		session.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Provider{})
	if err != nil {
		log.Println(err)
	}

	var providers []entity.Provider

	err1 := ps.BuildQuerySession().Find(&providers)
	if err1 != nil {
		log.Println(err1)
	}

	return NewPageData(total, providers)
}

func (self productService) FindAllAvailableProviders() (ps []entity.Provider) {
	_ = self.session.Where("enabled=?", true).Find(&ps)
	return
}

func (self productService) GetProviderById(id int64) (p entity.Provider, ok bool) {
	ok, _ = self.session.Where("id=?", id).Get(&p)
	return
}

func (self productService) GetProviderByProductId(id int64) (p entity.Provider, ok bool) {
	product, _ := self.GetProductById(id)
	ok, _ = self.session.Where("id=?", product.ProviderId).Get(&p)
	return
}

func (self productService) SaveProvider(p entity.Provider) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.session.Insert(&p)
		if err != nil {
			return
		}
		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetProviderById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.session.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此制造商不存在")
		}
	}
}

func (self productService) ToggleProviderEnabled(p *entity.Provider) error {
	p.Enabled = !p.Enabled
	_, err := self.session.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self productService) DeleteProvider(id int64) error {
	p, _ := self.GetProviderById(id)
	_, err := self.session.Delete(&p)
	return err
}

const detailFilePattern = "data/products/detail/%d.html"

func (self productService) SaveProductDetail(id int64, content string) (err error) {
	p, ok := self.GetProductById(id)
	if !ok {
		return errors.New("产品不存在！")
	}

	to := fmt.Sprintf(detailFilePattern, p.Id)
	_, err = os.Create(to)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(to, []byte(content), 0644)
	if err != nil {
		return
	}
	return
}

func (self productService) GetProductDetail(id int64) (detail string, err error) {
	from := fmt.Sprintf(detailFilePattern, id)
	f, err := os.Open(from)
	if err != nil {
		return
	}

	r, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	detail = string(r)
	return
}

func (self productService) GetProductImageFile(file string, t int) (targetFile *os.File, err error) {
	var dir string
	switch t {
	case entity.PTScheDiag:
		dir = revel.Config.StringDefault("dir.data.product.sd", "data/products/sd")
	case entity.PTPics:
		dir = revel.Config.StringDefault("dir.data.product.pics", "data/products/pics")
	default:
		return targetFile, errors.New("不支持类型")
	}
	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err = os.Open(imageFile)
	return
}

func (self productService) DeleteProductParams(id int64) (err error) {
	var p entity.ProductParams
	_, err = self.session.Where("id=?", id).Get(&p)
	if err != nil {
		return
	}
	_, err = self.session.Delete(&p)
	if err != nil {
		return
	}

	return nil
}

func (self productService) FindProductImages(id int64, t int) (ps []entity.ProductParams) {
	_ = self.session.Where("type=? and product_id=?", t, id).Find(&ps)
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Categories

func (self productService) FindAllCategoriesForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(session *xorm.Session) {
		session.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.ProductCategory{})
	if err != nil {
		log.Println(err)
	}

	var categories []entity.ProductCategory

	err1 := ps.BuildQuerySession().Find(&categories)
	if err1 != nil {
		log.Println(err1)
	}

	return NewPageData(total, categories)
}

func (self productService) availableQuery() *xorm.Session {
	return self.session.Where("enabled=?", true)
}

func (self productService) FindAllAvailableCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().Find(&ps)
	return
}

func (self productService) FindAvailableTopCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", 0).Find(&ps)
	return
}

func (self productService) FindAllAvailableCategoriesByParentId(id int64) (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", id).Find(&ps)
	return
}

func (self productService) GetCategoryById(id int64) (p entity.ProductCategory, ok bool) {
	ok, _ = self.session.Where("id=?", id).Get(&p)
	return
}

func (self productService) SaveCategory(p entity.ProductCategory) (id int64, err error) {
	log.Println("update")
	if p.Id == 0 { //insert
		_, err = self.session.Insert(&p)
		if err != nil {
			return
		}
		_ = self.updateCategoryCode(p)

		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetCategoryById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.session.Id(p.Id).Update(&p)
			log.Println("afeter update", p.ParentId, currDa.ParentId, p.DataVersion, currDa.DataVersion)
			if p.ParentId != currDa.ParentId {
				_ = self.updateCategoryCode(p)
			}
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此分类不存在")
		}
	}
}

func (self productService) updateCategoryCode(p entity.ProductCategory) error {
	p, _ = self.GetCategoryById(p.Id) //Hacked
	if p.ParentId != 0 {
		parent, _ := self.GetCategoryById(p.ParentId)
		p.Code = fmt.Sprintf("%v-%v", parent.Code, p.Id) //编码
	} else {
		p.Code = fmt.Sprintf("%v", p.Id) //编码
	}
	i, err := self.session.Id(p.Id).Cols("code").Update(&p)
	gotang.AssertNoError(err, "")
	log.Println("eff:", i)
	return err
}

func (self productService) ToggleCategoryEnabled(p *entity.ProductCategory) error {
	p.Enabled = !p.Enabled
	_, err := self.session.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self productService) FindAllProductStockLogs(id int64) (ps []entity.ProductStockLog) {
	_ = self.session.Where("product_id=?", id).Desc("id").Find(&ps)
	return
}

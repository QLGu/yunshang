package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

func NewProductService(db *xorm.Session) *ProductService {
	return &ProductService{db}
}

/////////////////////////////////////////////////////////////////////////////

type ProductService struct {
	db *xorm.Session
}

func (self ProductService) Total() int64 {
	total, _ := self.db.Count(&entity.Product{})
	return total
}

func (self ProductService) FindAllProductsForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Product{})
	gotang.AssertNoError(err, "")

	var products []entity.Product

	err1 := ps.BuildQuerySession().Find(&products)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, products)
}

func (self ProductService) FindAllAvailableProductsForPage(ps *PageSearcher) (page PageData) {
	ps.FilterCall = func(db *xorm.Session) {
		db.And("enabled=?", true)
	}

	ps.SearchKeyCall = func(db *xorm.Session) {
		db.And("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Product{})
	gotang.AssertNoError(err, "")

	var products []entity.Product
	err1 := ps.BuildQuerySession().Find(&products)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, products)
}

func (self ProductService) FindAllAvailableProducts() (ps []entity.Product) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
	return
}

func (self ProductService) GetProductById(id int64) (p entity.Product, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) SaveProduct(p entity.Product) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}
		p.Code = p.Id + entity.ProductStartDisplayCode //编码
		_, err = self.db.Id(p.Id).Cols("code").Update(&p)
		id = p.Id

		Emitter.Emit(EStockLog, "", p.Id, fmt.Sprintf("创建产品入库：%d(%s)", p.StockNumber, p.UnitName))

		return
	} else { // update
		currDa, ok := self.GetProductById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.db.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此产品不存在")
		}
	}
}

func (self ProductService) AddProductStock(id int64, stock int, message string) (newstock int, err error) {
	p, ok := self.GetProductById(id)
	if !ok {
		err = errors.New("产品不存在!")
		return
	}
	p.StockNumber += stock
	_, err = self.db.Id(p.Id).Cols("stock_number").Update(&p)
	Emitter.Emit(EStockLog, "", p.Id, fmt.Sprintf("入库：%d(%s), 当前库存%d(%s)", stock, p.UnitName, p.StockNumber, p.UnitName))

	newstock = p.StockNumber
	return
}

func (self ProductService) ToggleProductEnabled(p *entity.Product) error {
	p.Enabled = !p.Enabled
	if p.Enabled {
		p.EnabledAt = time.Now()
		_, err := self.db.Id(p.Id).Cols("enabled", "enabled_at").Update(p)
		return err
	}
	p.UnEnabledAt = time.Now()
	_, err := self.db.Id(p.Id).Cols("enabled", "un_enabled_at").Update(p)
	return err
}

func (self ProductService) FindAllProvidersForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Provider{})
	gotang.AssertNoError(err, "")

	var providers []entity.Provider
	err1 := ps.BuildQuerySession().Find(&providers)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, providers)
}

func (self ProductService) FindAllAvailableProviders() (ps []entity.Provider) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
	return
}

func (self ProductService) GetProviderById(id int64) (p entity.Provider, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) GetProviderByProductId(id int64) (p entity.Provider, ok bool) {
	product, _ := self.GetProductById(id)
	ok, _ = self.db.Where("id=?", product.ProviderId).Get(&p)
	return
}

func (self ProductService) SaveProvider(p entity.Provider) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}
		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetProviderById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.db.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此制造商不存在")
		}
	}
}

func (self ProductService) ToggleProviderEnabled(p *entity.Provider) error {
	p.Enabled = !p.Enabled
	_, err := self.db.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self ProductService) DeleteProvider(id int64) error {
	p, _ := self.GetProviderById(id)
	_, err := self.db.Delete(&p)
	return err
}

const detailFilePattern = "data/products/detail/%d.html"

func (self ProductService) SaveProductDetail(id int64, content string) (err error) {
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

func (self ProductService) GetProductDetail(id int64) (detail string, err error) {
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

func (self ProductService) GetProductImageFile(file string, t int) (targetFile *os.File, err error) {
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

func (self ProductService) DeleteProductParams(id int64) (err error) {
	var p entity.ProductParams
	_, err = self.db.Where("id=?", id).Get(&p)
	if err != nil {
		return
	}
	_, err = self.db.Delete(&p)
	if err != nil {
		return
	}

	return nil
}

func (self ProductService) FindProductImages(id int64, t int) (ps []entity.ProductParams) {
	_ = self.db.Where("type=? and product_id=?", t, id).Find(&ps)
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Categories

func (self ProductService) FindAllCategoriesForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.ProductCategory{})
	gotang.AssertNoError(err, "")

	var categories []entity.ProductCategory
	err1 := ps.BuildQuerySession().Find(&categories)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, categories)
}

func (self ProductService) availableQuery() *xorm.Session {
	return self.db.Where("enabled=?", true)
}

func (self ProductService) FindAllAvailableCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().Find(&ps)
	return
}

func (self ProductService) FindAvailableTopCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", 0).Find(&ps)
	return
}

func (self ProductService) FindAllAvailableCategoriesByParentId(id int64) (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", id).Find(&ps)
	return
}

func (self ProductService) GetCategoryById(id int64) (p entity.ProductCategory, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) SaveCategory(p entity.ProductCategory) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
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
			_, err = self.db.Id(p.Id).Update(&p)
			if p.ParentId != currDa.ParentId {
				_ = self.updateCategoryCode(p)
			}
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此分类不存在")
		}
	}
}

func (self ProductService) updateCategoryCode(p entity.ProductCategory) (err error) {
	p, _ = self.GetCategoryById(p.Id) //Hacked
	if p.ParentId != 0 {
		parent, _ := self.GetCategoryById(p.ParentId)
		p.Code = fmt.Sprintf("%v-%v", parent.Code, p.Id) //编码
	} else {
		p.Code = fmt.Sprintf("%v", p.Id) //编码
	}
	_, err = self.db.Id(p.Id).Cols("code").Update(&p)
	return
}

func (self ProductService) ToggleCategoryEnabled(p *entity.ProductCategory) error {
	p.Enabled = !p.Enabled
	_, err := self.db.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self ProductService) FindAllProductStockLogs(id int64) (ps []entity.ProductStockLog) {
	_ = self.db.Where("product_id=?", id).Desc("id").Find(&ps)
	return
}

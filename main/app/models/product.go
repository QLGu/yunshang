package models

import (
	"fmt"
	"log"
	"time"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

type ProductService interface {
	Total() int64

	FindAllAvailableProducts() []entity.Product

	FindAllProductsForPage(ps *PageSearcher) PageData

	SaveProduct(p entity.Product) (id int64, err error)

	GetProductById(id int64) (entity.Product, bool)

	ToggleProductEnabled(p *entity.Product) error

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

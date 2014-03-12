package models

import (
	"fmt"
	"log"
	"time"

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

func (self productService) FindAllAvailableProducts() (ps []entity.Product) {
	_ = self.session.Where("enabled=?", true).Find(&ps)
	return
}
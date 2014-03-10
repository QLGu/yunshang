package models

import (
	"log"

	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

type ProductService interface {
	Total() int64

	FindAllProductsForPage(ps *PageSearcher) PageData
}

func NewProductService(session *xorm.Session) ProductService {
	return &productService{}
}

/////////////////////////////////////////////////////////////////////////////

type productService struct {
}

func (self productService) Total() int64 {
	return 0
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

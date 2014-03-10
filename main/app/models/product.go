package models

import (
	"github.com/lunny/xorm"
)

type ProductService interface {
	Total() int64
}

func NewProductService(session *xorm.Session) ProductService {
	return &productService{}
}

/////////////////////////////////////////////////////////////////////////////

type productService struct {
}

func (e productService) Total() int64 {
	return 0
}

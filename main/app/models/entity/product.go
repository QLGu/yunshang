package entity

import (
	"time"
)

const ProductStartDisplayCode int64 = 10000

// 产品
type Product struct {
	Id int64 `json:"id"`

	Code      int64  `xorm:"unique" json:"code"` //商品编号
	Name      string `json:"name"`               // 名称
	NameExtra string `json:"name_extra`          // 附加名称
	Category  string `json:"category"`           // 分类

	Model string `json:"model"` //型号

	UnitName    string `json:"unit_name"`    // 商品计量单位
	StockNumber int    `json:"stock_number"` //库存数量

	Provider  string `json:"provider"`                      // 制造商/供应商
	Introduce string `xorm:"varchar(1000)" json:"introduce` //简介

	ScoresLevel int `json:"scores_level` // 评价等级

	Enabled bool `json:"enabled"` // 上架/下架

	EnabledAt   time.Time `json:"enabled_at"`   //上架时间
	UnEnabledAt time.Time `json:"unenabled_at"` //下架时间

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`

	//TODO images
	// 说明页
	// 资料文档
	// 服务指南
}

// 产品定价规则
type ProductPriceRule struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`

	Name          string  `json:"name"`           // 名称
	PriceType     string  `json:"name"`           // 一口价|按量计价
	Price         float64 `json:"price"`          // 单价
	StartQuantity int     `json:"start_quantity"` // 按量计价起始数量
	EndQuantity   int     `json:"end_quantity"`   // 按量计价最大数量， 0 表示 无限制
}

// 产品详细参数
type ProductSpec struct {
	Id        int64  `json:"id"`
	ProductId int64  `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

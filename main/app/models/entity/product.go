package entity

import (
	"time"
)

const ProductStartDisplayCode int64 = 10000

type ProductCategory struct {
	Id      int64 `json:"id"`
	Enabled bool  `json:"enabled"`

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`

	ParentId    int64  `json:"parent_id"`   //上一级ID
	Name        string `json:"name"`        //名称
	Code        string `json:"code"`        //编码
	Description string `json:"description"` //描述

	Tags string `json:"tags"` // 标签
}

// 产品
type Product struct {
	Id int64 `json:"id"`

	Code         int64  `xorm:"unique" json:"code"` //商品编号
	Name         string `json:"name"`               // 名称
	NameExtra    string `json:"name_extra`          // 附加名称
	CategoryId   int64  `json:"category_id"`        // 分类ID
	CategoryCode string `json:"category_code`       // 产品分类code

	Model string `json:"model"` //型号

	UnitName          string `json:"unit_name"`            // 商品计量单位
	StockNumber       int    `json:"stock_number"`         //库存数量
	MinNumberOfOrders int    `json:"min_number_of_orders"` //最小订货数量

	Price float64 `xorm:"Numeric" json:"price"` //单价， 冗余字段

	ProviderId int64  `json:"provider"`                       // 制造商/供应商Id
	Introduce  string `xorm:"varchar(1000)" json:"introduce"` //简介

	LotNumber string `json:"lot_number"` //批号
	Delivery  string `json:"delivery"`   //货期

	ScoresLevel int `json:"scores_level` // 评价等级

	Enabled bool `json:"enabled"` // 上架/下架

	Tags      string `json:"tags"`       // 标签
	SortValue int    `json:"sort_value"` //排序值

	EnabledAt   time.Time `json:"enabled_at"`   //上架时间
	UnEnabledAt time.Time `json:"unenabled_at"` //下架时间

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`

	// images
	// 说明页
	// 资料文档
	// 服务指南
}

func (e Product) DisplayName() string {
	if len(e.Model) > 0 {
		return e.Model
	}
	return e.Name
}

// 产品定价规则
type ProductPrices struct {
	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`

	Name          string  `json:"name"`                 // 名称
	Price         float64 `xorm:"Numeric" json:"price"` // 单价
	StartQuantity int     `json:"start_quantity"`       // 按量计价起始数量
	EndQuantity   int     `json:"end_quantity"`         // 按量计价最大数量， 0 表示 无限制
}

const (
	PTSpec     = 1 // 参数
	PTScheDiag = 2 //示意图
	PTPics     = 3 // 图库
	PTMaterial = 4 // 资料
)

// 产品详细参数
type ProductParams struct {
	Id        int64  `json:"id"`
	ProductId int64  `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`

	Type int `json:"type"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

// 产品入库记录
type ProductStockLog struct {
	Id int64 `json:"id"`

	ProductId int64     `json:"product_id"`
	Message   string    `json:"message"`
	User      string    `json:"user"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

// 产品制造商
type Provider struct {
	Id      int64 `json:"id"`
	Enabled bool  `json:"enabled"`

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`

	Name       string `json:"name"`                          // 名称
	ShortName  string `json:"short_name"`                    // 简称
	Introduce  string `xorm:"varchar(1000)" json:"introduce` //简介
	MainBiz    string `xorm:"varchar(1000)" json:"main_biz"` // 主要产品或服务
	WebsiteUrl string `json:"website_url"`                   // 公司主页
	Tags       string `json:"tags"`                          // 标签
}

func (e Provider) DisplayName() string {
	if len(e.ShortName) > 0 {
		return e.ShortName
	}
	return e.Name
}

// 产品收藏记录
type ProductCollect struct {
	Id int64 `json:"id"`

	ProductId int64     `json:"product_id"`
	UserId    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

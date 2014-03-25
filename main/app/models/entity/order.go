package entity

import (
	"time"
)

// 购物车
type Cart struct {
	Id int64 `json:"id"`

	UserId int64 `json:"user_id"`

	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

type Payment struct {
	Id int64 `json:"id"`

	Name        string `json:"name"`
	Description string `xorm:"varchar(1000)" json:"description"`

	Enabled bool `json:"enabled"`
}

const (
	OS_TEMP   = 1 //临时订单
	OS_SUBMIT = 2 //提交的订单
	OS_PAY    = 3 // 支付的订单
	OS_CANEL  = 4 // 取消的订单
	OS_LOCK   = 5 // 锁定的订单
)

// 订单
type Order struct {
	Id   int64 `json:"id"`
	Code int64 `json:"code"` //订单号

	UserId int64 `json:"user_id"`

	DaId      int64 `json:"da_id"`      // 收货地址
	PaymentId int64 `json:"payment_id"` //支付方式
	InvoiceId int64 `json:"invoice_id"` //发票信息

	Amount float64 `xorm:"Numeric" json:"amount"` //总计

	SubmitAt time.Time `json:"submit_at"` //提交时间

	Status int `json:"status"` //状态

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

//订单明细
type OrderDetail struct {
	Id        int64   `json:"id"`
	OrderId   int64   `json:"order_id"`
	ProductId int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `xorm:"Numeric" json:"price"` //单价
}

package entity

import (
	"time"
)

type Cart struct {
	Id int64 `json:"id"`

	UserId int64 `json:"user_id"`

	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

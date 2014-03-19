package entity

import (
	"time"
)

const (
	ATAd = 1 //
)

type AppParams struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`
	Value string `json:"value"`

	Type int `json:"type"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

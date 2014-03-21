package entity

import (
	"time"
)

const (
	ATAd = 1 // 广告图片
	ATHk = 2 // 关键词
	ATSg = 3 //标语
)

type AppParams struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`
	Value string `json:"value"`

	Type int `json:"type"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

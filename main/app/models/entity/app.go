package entity

import (
	"time"
)

const (
	ATAd = 1 // 广告图片
	ATHk = 2 // 关键词
	ATSg = 3 // 标语
)

//应用参数
type AppParams struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`
	Value string `json:"value"`

	Type int `json:"type"`

	Data string `json:"data"` //扩展信息

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

//询价
type Inquiry struct {
	Id     int64 `json:"id"`
	UserId int64 `json:"user_id"`

	Model    string `json:"model"`
	Quantity int    `json:"quantity"`
	Contact  string `json:"contact"`
	Phone    string `json:"phone"`
	QqMsn    string `json:"qq_msn"`

	Replies int `json:"replies"` //回复数量， 冗余字段

	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

//询价回复
type InquiryReply struct {
	Id        int64 `json:"id"`
	InquiryId int64 `json:"inquiry_id"`
	UserId    int64 `json:"user_id"`

	Title string `json:"title"`

	Content string `xorm:"varchar(1000)" json:"content"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

//Migrations
type Migration struct {
	Id   int64  `json:"id"`
	Name string `xorm:"unique" json:"name"`

	Description string    `json:"description"`
	CreatedAt   time.Time `xorm:"created" json:"created_at"`
}

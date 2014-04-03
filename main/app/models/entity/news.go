package entity

import (
	"strings"
	"time"
)

const NewsStartDisplayCode = 10000

// 新闻分类
type NewsCategory struct {
	Id      int64 `json:"id"`
	Enabled bool  `json:"enabled"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`

	ParentId    int64  `json:"parent_id"`   //上一级ID
	Name        string `json:"name"`        //名称
	Code        string `json:"code"`        //编码
	Description string `json:"description"` //描述
}

// 新闻
type News struct {
	Id   int64 `json:"id"`
	Code int64 `xorm:"unique" json:"code"` //新闻编号

	CategoryId   int64  `json:"category_id"`  // 分类ID
	CategoryCode string `json:"category_code` // 分类code

	UserId   int64  `json:"user_id"`   //发布人
	UserName string `json:"user_name"` //发布人姓名

	Hits      int64     `json:"hits"`    //点击量
	Enabled   bool      `json:"enabled"` //是否发布?
	Source    string    `json:"source"`  //来源
	PublishAt time.Time `json:"publish_at"`

	Title    string `json:"title"`                        //标题
	Subtitle string `json:"subtitle"`                     //副标题
	Summary  string `xorm:"varchar(1000)" json:"summary"` //摘要
	Content  string `json:"content"`                      //内容

	Tags string `json:"tags"` //标签

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

//是客户服务的文章？
// 硬编码
func (e News) IsServiceArticle() bool {
	return strings.HasPrefix(e.CategoryCode, "4")
}

func (e News) IsAboutArticle() bool {
	return strings.HasPrefix(e.CategoryCode, "10")
}

func (e News) IsPureNews() bool {
	return !e.IsServiceArticle() && !e.IsAboutArticle()
}

func (e News) DisplayAt() time.Time {
	if e.PublishAt.IsZero() {
		if !e.UpdatedAt.IsZero() {
			return e.UpdatedAt
		}
		return e.CreatedAt
	}

	return e.PublishAt
}

const (
	NTScheDiag = 1 //示意图
	NTPics     = 2 // 图库
	NTMaterial = 3 // 资料
)

// 新闻参数
type NewsParam struct {
	Id     int64  `json:"id"`
	NewsId int64  `json:"news_id"`
	Name   string `json:"name"`
	Value  string `json:"value"`

	Type int `json:"type"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

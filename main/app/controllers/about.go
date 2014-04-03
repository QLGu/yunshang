package controllers

import (
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/revel/revel"
)

// 相关Actions
type About struct {
	AppController
}

// 主页
func (c About) Index() revel.Result {
	return c.Redirect(routes.About.View(0))
}

// 主页
func (c About) View(id int64) revel.Result {
	newslist := c.newsApi().FindAllAvailableNewsByCategory(10)

	var current entity.News
	if id == 0 {
		current = newslist[0]
	} else {
		var exists = true
		current, exists = c.newsApi().GetNewsById(id)
		if !exists {
			return c.NotFound("此条目不存在！")
		}
	}
	currentDetail, err := c.newsApi().GetNewsDetail(current.Id)
	if err != nil {
		currentDetail = ""
	}

	c.setChannel("about/view")
	return c.Render(newslist, current, currentDetail)
}

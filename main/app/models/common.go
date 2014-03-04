package models

import (
	"reflect"
	"regexp"

	"github.com/itang/gotang"
	"github.com/lunny/xorm"
)

// 会话回调处理
type PageSearcherCall func(session *xorm.Session)

// 分页搜索器
type PageSearcher struct {
	Limit         int
	Start         int
	SortField     string
	SortDir       string
	FilterCall    PageSearcherCall
	SearchKeyCall PageSearcherCall
	OtherCalls    []PageSearcherCall
	Search        string
	Session       *xorm.Session
}

// 分页数据
type PageData struct {
	Total int64
	Data  interface{}
}

// 构造新的分页数据
func NewPageData(total int64, data interface{}) PageData {
	if reflect.ValueOf(data).IsNil() {
		return PageData{total, []string{}}
	}
	return PageData{total, data}
}

// 构建Count会话
func (self *PageSearcher) BuildCountSession() *xorm.Session {
	self.doCommon()

	return self.Session
}

// 构建查询会话
func (self *PageSearcher) BuildQuerySession() *xorm.Session {
	self.doCommon()

	if len(self.SortField) != 0 {
		if self.SortDir == "desc" {
			self.Session.Desc(self.SortField)
		} else {
			self.Session.Asc(self.SortField)
		}
	} else {
		self.Session.Desc("id")
	}

	self.Session.Limit(self.Limit, self.Start)

	if len(self.OtherCalls) != 0 {
		for _, call := range self.OtherCalls {
			call(self.Session)
		}
	}

	return self.Session
}

// 执行常用处理
func (self *PageSearcher) doCommon() {
	gotang.Assert(self.Session != nil, "设置session先")

	re := regexp.MustCompile(".*([';]+|(--)+).*")
	self.Search = re.ReplaceAllString(self.Search, " ")

	if self.SearchKeyCall != nil && len(self.Search) != 0 {
		self.SearchKeyCall(self.Session)
	}

	if self.FilterCall != nil {
		self.FilterCall(self.Session)
	}
}

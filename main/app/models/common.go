package models

import (
	"fmt"
	"reflect"
	"regexp"

	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/itang/gotang"
	"github.com/lunny/xorm"
)

//业务错误/BizError/可恢复错误
type BizError struct {
	error
}

func NewBizError(err error) BizError {
	return BizError{err}
}

// 分页数据
type PageData struct {
	Total int64
	Data  interface{}

	page  int64
	limit int64
}

func (e PageData) IsLastPage() bool {
	return e.page == e.Pages()
}

func (e PageData) IsFirstPage() bool {
	return e.page == 1
}

func (e PageData) HasNextPage() bool {
	return e.page < e.Pages()
}

func (e PageData) NextPage() int64 {
	return e.page + 1
}

func (e PageData) PrevPage() int64 {
	return e.page - 1
}

func (e PageData) HasPrevPage() bool {
	return e.page > 1
}

func (e PageData) Pages() int64 {
	return (e.Total-1)/e.limit + 1
}

func (e PageData) Page() int64 {
	return e.page
}

func (e PageData) Limit() int64 {
	return e.limit
}

func (e PageData) Start() int64 {
	return e.limit * (e.page - 1)
}

func (e *PageData) SetPage(page int64) {
	if page <= 0 {
		e.page = 1
	} else if page > e.Pages() {
		e.page = e.Pages()
	} else {
		e.page = page
	}
}

func (e *PageData) SetLimit(limit int64) {
	if limit <= 0 {
		e.limit = 10
	} else {
		e.limit = limit
	}
	gotang.Assert(e.limit != 0, "limit should not be zero!")
}

func (e *PageData) PageNumbers() []int {
	ret := make([]int, 0)
	pages := int(e.Pages())
	if pages == 1 {
		return []int{1}
	}
	for i := 1; i <= pages; i++ {
		ret = append(ret, i)
	}
	return ret
}

// 会话回调处理
type PageSearcherCall func(db *xorm.Session)

// 分页搜索器
type PageSearcher struct {
	Limit         int64
	Start         int64
	Page          int64
	SortField     string
	SortDir       string
	FilterCall    PageSearcherCall
	SearchKeyCall PageSearcherCall
	OtherCalls    []PageSearcherCall
	Search        string
	db            *xorm.Session
}

func (e *PageSearcher) SetDb(db *xorm.Session) {
	e.db = db
}

// 构造新的分页数据
func NewPageData(total int64, data interface{}, ps *PageSearcher) *PageData {
	//start : from 0
	var page *PageData
	if reflect.ValueOf(data).IsNil() {
		page = &PageData{Total: total, Data: []string{}}
	} else {
		page = &PageData{Total: total, Data: data}
	}
	page.SetLimit(ps.Limit)
	page.SetPage(ps.Page)
	return page
}

// 构建Count会话
func (self *PageSearcher) BuildCountSession() *xorm.Session {
	self.doCommon()

	return self.db
}

// 构建查询会话
func (self *PageSearcher) BuildQuerySession() *xorm.Session {
	self.doCommon()

	if len(self.SortField) != 0 {
		if self.SortDir == "desc" {
			self.db.Desc(self.SortField)
		} else {
			self.db.Asc(self.SortField)
		}
	} else {
		self.db.Desc("id")
	}

	self.db.Limit(int(self.Limit), int(self.Start))

	return self.db
}

// 执行常用处理
func (self *PageSearcher) doCommon() {
	gotang.Assert(self.db != nil, "设置db先")

	re := regexp.MustCompile(".*([';]+|(--)+).*")
	self.Search = re.ReplaceAllString(self.Search, " ")

	if self.SearchKeyCall != nil && len(self.Search) != 0 {
		self.SearchKeyCall(self.db)
	}

	if self.FilterCall != nil {
		self.FilterCall(self.db)
	}

	if len(self.OtherCalls) != 0 {
		for _, call := range self.OtherCalls {
			call(self.db)
		}
	}
}

// go-linq helper
func eqSlice(s1 []int, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func asStrSlice(is []int) []string {
	ss := make([]string, len(is))
	for i, v := range is {
		ss[i] = fmt.Sprintf("%d", v)
	}
	return ss
}

func asStrSliceFromInt64(is []int64) []string {
	ss := make([]string, len(is))
	for i, v := range is {
		ss[i] = fmt.Sprintf("%d", v)
	}
	return ss
}

func asIntSlice(ts []T) []int {
	is := make([]int, len(ts))
	for i, v := range ts {
		is[i] = v.(int)
	}
	return is
}

func asInt64Slice(ts []T) []int64 {
	is := make([]int64, len(ts))
	for i, v := range ts {
		is[i] = v.(int64)
	}
	return is
}

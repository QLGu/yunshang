package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	gio "github.com/itang/gotang/io"
	//. "github.com/ahmetalpbalkan/go-linq"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
	"strconv"
)

////////////////////////////////////////////////////////
type NewsService struct {
	db *xorm.Session
}

func NewNewsService(db *xorm.Session) *NewsService {
	return &NewsService{db}
}

////////////////////////////////////
func (self NewsService) FindAllNewsForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("title like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.News{})
	gotang.AssertNoError(err, "")

	var news []entity.News

	err1 := ps.BuildQuerySession().Find(&news)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, news, ps)
}

func (self NewsService) SaveNews(p entity.News) (id int64, err error) {
	if p.Id == 0 { //insert
		p.CategoryCode = self.GetCategoryCode(p.CategoryId)
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}

		p.Code = p.Id + entity.NewsStartDisplayCode //编码
		_, err = self.db.Id(p.Id).Cols("code").Update(&p)
		id = p.Id

		return
	} else { // update
		currDa, ok := self.GetNewsById(p.Id)
		if ok {
			if p.CategoryId != currDa.CategoryId || len(currDa.CategoryCode) == 0 { //检测类型是否修改
				p.CategoryCode = self.GetCategoryCode(p.CategoryId) //更新分类编码， 冗余设计
			}

			_, err = self.db.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此新闻不存在")
		}
	}
}

func (self NewsService) ToggleNewsEnabled(p *entity.News) error {
	p.Enabled = !p.Enabled
	if p.Enabled {
		p.PublishAt = time.Now()
		_, err := self.db.Id(p.Id).Cols("enabled", "publish_at").Update(p)
		return err
	}
	p.PublishAt = time.Time{}
	_, err := self.db.Id(p.Id).Cols("enabled", "publish_at").Update(p)
	return err
}

func (self NewsService) GetNewsById(id int64) (n entity.News, exists bool) {
	exists, _ = self.db.Where("id=?", id).Get(&n)
	return
}

const newsDetailFilePattern = "data/news/detail/%d.html"

func (self NewsService) SaveNewsDetail(id int64, content string) (err error) {
	p, ok := self.GetNewsById(id)
	if !ok {
		return errors.New("新闻不存在！")
	}

	to := fmt.Sprintf(newsDetailFilePattern, p.Id)
	_, err = os.Create(to)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(to, []byte(content), 0644)
	if err != nil {
		return
	}
	return
}

func (self NewsService) GetNewsDetail(id int64) (detail string, err error) {
	from := fmt.Sprintf(newsDetailFilePattern, id)
	f, err := os.Open(from)
	if err != nil {
		return
	}

	r, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	detail = string(r)
	return
}

func (self NewsService) FindAllAvailableNewsByCategory(ctId int64) (ps []entity.News) {
	_ = self.db.Where("enabled=? and category_id=?", true, ctId).Asc("id").Find(&ps)
	return
}

////////////////////////////////////
func (self NewsService) FindAllAvailableCategories() (ps []entity.NewsCategory) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
	return
}

func (self NewsService) FindAllAvailableServiceCategories() (ps []entity.NewsCategory) {
	_ = self.db.Where("enabled=? and parent_id=?", true, 4).Asc("id").Find(&ps)
	return
}

func (self NewsService) GetCategoryCode(id int64) string {
	c, _ := self.GetCategoryById(id)
	return c.Code
}

func (self NewsService) GetCategoryById(id int64) (e entity.NewsCategory, exists bool) {
	exists, _ = self.db.Where("id=?", id).Get(&e)
	return
}

func (self NewsService) SaveCategory(p entity.NewsCategory) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}
		_ = self.updateCategoryCode(p)

		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetCategoryById(p.Id)
		if ok {
			_, err = self.db.Id(p.Id).Update(&p)
			if p.ParentId != currDa.ParentId {
				_ = self.updateCategoryCode(p)
			}
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此分类不存在")
		}
	}
}

func (self NewsService) updateCategoryCode(p entity.NewsCategory) (err error) {
	p, _ = self.GetCategoryById(p.Id) //Hacked
	if p.ParentId != 0 {
		parent, _ := self.GetCategoryById(p.ParentId)
		p.Code = fmt.Sprintf("%v-%v", parent.Code, p.Id) //编码
	} else {
		p.Code = fmt.Sprintf("%v", p.Id) //编码
	}
	_, err = self.db.Id(p.Id).Cols("code").Update(&p)
	return
}

func (self NewsService) FindAllCategoriesForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.NewsCategory{})
	gotang.AssertNoError(err, "")

	var ncs []entity.NewsCategory

	err1 := ps.BuildQuerySession().Find(&ncs)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, ncs, ps)
}

func (self NewsService) ToggleCategoryEnabled(p *entity.NewsCategory) error {
	p.Enabled = !p.Enabled
	_, err := self.db.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self NewsService) NewsCommentsForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("content like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Comment{})
	gotang.AssertNoError(err, "")

	var comments []entity.Comment
	err1 := ps.BuildQuerySession().Find(&comments)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, comments, ps)
}

func (self NewsService) GetNewsImageFile(file string, t int) (targetFile *os.File, err error) {
	var dir string
	switch t {
	case entity.NTScheDiag:
		dir = revel.Config.StringDefault("dir.data.news.sd", "data/news/sd")
	case entity.NTPics:
		dir = revel.Config.StringDefault("dir.data.news.pics", "data/news/pics")
	default:
		return targetFile, errors.New("不支持类型")
	}
	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err = os.Open(imageFile)
	return
}

func (self NewsService) DeleteNewsParam(id int64) (err error) {
	var p entity.NewsParam
	_, err = self.db.Where("id=?", id).Get(&p)
	if err != nil {
		return
	}
	_, err = self.db.Delete(&p)
	if err != nil {
		return
	}

	return nil
}

func (self NewsService) FindNewsImages(id int64, t int) (ps []entity.NewsParam) {
	_ = self.db.Where("type=? and news_id=?", t, id).Find(&ps)
	return
}

func (self NewsService) FindNewsMaterial(id int64) (files []entity.NewsParam) {
	_ = self.db.Where("type=? and news_id=?", entity.NTMaterial, id).Find(&files)
	return
}

func (self NewsService) FindTWNews(ctcode string, limit int) (ps []entity.News) {
	_ = self.AvailableWithCtCodeLimit(ctcode, limit).And("tags like ?", "%#图文%").Find(&ps)
	return
}

func (self NewsService) FindNoTWNews(ctcode string, limit int) (ps []entity.News) {
	_ = self.AvailableWithCtCodeLimit(ctcode, limit).And("tags not like ?", "%#图文%").Find(&ps)
	return
}

func (self NewsService) FindNews(ctcode string, limit int) (ps []entity.News) {
	_ = self.AvailableWithCtCodeLimit(ctcode, limit).Find(&ps)
	return
}

func (self NewsService) FindPureNews(ctcode string, limit int) (ps []entity.News) {
	_ = self.AvailableWithCtCodeLimit(ctcode, limit).And("(category_code not like ?  and category_code not like ?)", "10%", "4%").Find(&ps)
	return
}

func (self NewsService) Available() *xorm.Session {
	return self.db.And("enabled=true")
}

func (self NewsService) AvailableWithCtCodeLimit(ctcode string, limit int) (session *xorm.Session) {
	session = self.Available().Limit(limit).Desc("id")
	if len(ctcode) > 0 {
		session.And("(category_code =? or category_code like ?)", ctcode, ctcode+"-%")
	}
	return
}

func (self NewsService) GetCategoryByCodeString(ctcode string) (p entity.NewsCategory, exists bool) {
	codestr := ""
	cArr := strings.Split(ctcode, "-")
	if len(cArr) > 0 {
		codestr = cArr[len(cArr)-1] //最后一个
	}
	code, err := strconv.Atoi(codestr)
	if err != nil {
		return
	}
	exists, _ = self.db.Where("code=?", code).Get(&p)
	return
}

func (self NewsService) GetNextDisplayNews(currid int64) []entity.DisplayItem {
	rows, err := self.db.Query("select id, title from t_news where id = (select min(id) from t_news where id > ? and enabled=true)", currid)
	gotang.AssertNoError(err, "")

	return self.rowsAsDisplayItems(rows)
}

func (self NewsService) GetPrevDisplayNews(currid int64) []entity.DisplayItem {
	rows, err := self.db.Query("select id, title from t_news where id = (select max(id) from t_news where id < ? and enabled=true)", currid)
	gotang.AssertNoError(err, "")

	return self.rowsAsDisplayItems(rows)
}

func (self NewsService) rowsAsDisplayItems(rows []map[string][]byte) (ns []entity.DisplayItem) {
	if len(rows) == 0 {
		return
	}

	row := rows[0]
	idstr := string(row["id"])
	title := string(row["title"])
	id, err := strconv.Atoi(idstr)
	gotang.AssertNoError(err, "id is integer")

	return []entity.DisplayItem{{int64(id), title}}
}

func (self NewsService) TotalNewsComments() int64 {
	total, err := self.db.Where("target_type=?", entity.CT_ARTICLE).Count(&entity.Comment{})
	gotang.AssertNoError(err, "TotalNewsComments")
	return total
}

func (self NewsService) TotalNewsCommentsUnconfirm() int64 {
	total, err := self.db.Where("target_type=? and enabled=false", entity.CT_ARTICLE).Count(&entity.Comment{})
	gotang.AssertNoError(err, "TotalNewsCommentsUnconfirm")
	return total
}

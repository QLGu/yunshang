package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	gio "github.com/itang/gotang/io"
	//. "github.com/ahmetalpbalkan/go-linq"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
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

////////////////////////////////////
func (self NewsService) FindAllAvailableCategories() (ps []entity.NewsCategory) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
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

func (self NewsService) SaveCategory(e entity.NewsCategory) (newid int64, err error) {
	_, err = self.db.Insert(&e)
	newid = e.Id

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

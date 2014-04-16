package models

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

func NewAppService(session *xorm.Session) *AppService {
	return &AppService{session}
}

type AppService struct {
	db *xorm.Session
}

func noFoundError(id int64, t string) BizError {
	return NewBizError(fmt.Errorf("Id为%d的%s不存在", id, t))
}

func (self AppService) FindAdImages() (ps []entity.AppParams) {
	err := self.db.Where("type=?", entity.ATAd).Desc("updated_at").Find(&ps)
	AssertNoError(err, "FindAdImages")

	return
}

func (self AppService) ForceGetAppParamsById(id int64) (e entity.AppParams) {
	has, err := self.db.Id(id).Get(&e)
	AssertNoError(err, "ForceGetAppParamsById")
	Assert(has, fmt.Sprintf("Id为%d的%s不存在", id, "AppParams"))

	return
}

func (self AppService) GetAdImageFile(file string) (targetFile *os.File, err error) {
	var dir string = revel.Config.StringDefault("dir.data.adimages", "data/adimages")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err = os.Open(imageFile)
	return
}

func (self AppService) SetFirstAdImage(id int64) (bizErr BizError) {
	var ad = self.ForceGetAppParamsById(id)

	affected, err := self.db.Id(ad.Id).Update(&ad)
	AssertNoError(err, "SetFirstAdImage, Update")
	Assert(affected == 1, "删除未成功")

	return
}

func (self AppService) SetAdImageLink(id int64, link string) (bizErr BizError) {
	var ad = self.ForceGetAppParamsById(id)
	ad.Data = link
	affected, err := self.db.Id(ad.Id).Cols("data").Update(&ad)
	AssertNoError(err, "SetAdImageLink, Update")
	Assert(affected == 1, "更新未成功")

	return
}

func (self AppService) DeleteAdImage(id int64) (bizErr BizError) {
	var ad = self.ForceGetAppParamsById(id)

	affected, err := self.db.Delete(&ad)
	AssertNoError(err, "DeleteAdImage, Delete")
	Assert(affected == 1, "删除未成功")

	return
}

func (self AppService) FindHotKeywords() (ps []entity.AppParams) {
	err := self.db.Where("type=?", entity.ATHk).Desc("updated_at").Find(&ps)
	AssertNoError(err, "FindHotKeywords")

	return
}

func (self AppService) DeleteHotKeyword(id int64) (err error) {
	var ad = self.ForceGetAppParamsById(id)

	affected, err := self.db.Delete(&ad)
	AssertNoError(err, "DeleteHotKeyword, Delete")
	Assert(affected == 1, "删除未成功")

	return err
}

func (self AppService) SetFirstHotKeyword(id int64) (bizErr BizError) {
	var ad = self.ForceGetAppParamsById(id)

	affected, err := self.db.Id(ad.Id).Update(&ad)
	AssertNoError(err, "SetFirstHotKeyword, Update")
	Assert(affected == 1, "更新未成功")

	return
}

func (self AppService) GetSloganContent() string {
	p, exists := self.GetSlogan()
	if !exists {
		return ""
	}
	return p.Value
}

func (self AppService) GetSlogan() (s entity.AppParams, ok bool) {
	var sgs []entity.AppParams
	_ = self.db.Where("type=?", entity.ATSg).Find(&sgs)
	if len(sgs) > 0 {
		return sgs[0], true
	}
	return s, false
}

func (self AppService) SaveSlogan(s entity.AppParams) (bizError BizError) {
	if s.Id == 0 { //new
		s.Type = entity.ATSg
		affected, err := self.db.Insert(&s)
		AssertNoError(err, "SaveSlogan, insert")
		Assert(affected == 1, "插入未成功")
	} else {
		p, exists := self.GetSlogan()
		if exists {
			p.Value = s.Value
			affected, err := self.db.Id(p.Id).Cols("value").Update(&p)
			AssertNoError(err, "SaveSlogan, err")
			Assert(affected == 1, "更新未成功")
		}
	}
	return
}

func (self AppService) SavePayment(ps ...entity.Payment) (err error) {
	affected, err := self.db.Insert(&ps)
	Assert(affected > 0, "保存未成功")
	return
}

func (self AppService) SaveInquiry(i entity.Inquiry) (err error) {
	_, err = self.db.Insert(&i)
	return
}

func (self AppService) FindAllMigrations() (ps []entity.Migration) {
	_ = self.db.Find(&ps)
	return
}

func (self AppService) FindAllMigrationsAsMap() map[string]bool {
	ps := self.FindAllMigrations()
	ret := make(map[string]bool)
	for _, p := range ps {
		ret[p.Name] = true
	}
	return ret
}

func (self AppService) ExistsMigrations() bool {
	c, err := self.db.Count(&entity.Migration{})
	AssertNoError(err, "ExistsMigrations")

	return c > 0
}

func (self AppService) SaveMigration(name string, desc string) (err error) {
	_, err = self.db.Insert(entity.Migration{Name: name, Description: desc})
	AssertNoError(err, "SaveMigration")

	return
}

func (self AppService) FindAllFeedbacksForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		lv := "%" + ps.Search + "%"
		db.Where("subject like ? or name like ? or content like ?", lv, lv, lv)
	}

	total, err := ps.BuildCountSession().Count(&entity.Feedback{})
	AssertNoError(err, "")

	var ins []entity.Feedback
	err1 := ps.BuildQuerySession().Find(&ins)
	AssertNoError(err1, "")

	return NewPageData(total, ins, ps)
}

func (self AppService) TotalFeedbacks() int64 {
	total, err := self.db.Count(&entity.Feedback{})
	AssertNoError(err, "TotalFeedbacks")

	return total
}

func (self AppService) SaveFeedback(p entity.Feedback) (err error) {
	_, err = self.db.Insert(&p)
	return
}

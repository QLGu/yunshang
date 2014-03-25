package models

import (
	"os"
	"path/filepath"

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

func (self AppService) FindAdImages() (ps []entity.AppParams) {
	_ = self.db.Where("type=?", entity.ATAd).Desc("updated_at").Find(&ps)
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

func (self AppService) SetFirstAdImage(id int64) (err error) {
	var ad entity.AppParams
	_, err = self.db.Id(id).Get(&ad)
	if err != nil {
		return
	}
	_, err = self.db.Id(ad.Id).Update(&ad)
	return err
}

func (self AppService) DeleteAdImage(id int64) (err error) {
	var ad entity.AppParams
	_, err = self.db.Id(id).Get(&ad)
	if err != nil {
		return
	}
	_, err = self.db.Delete(&ad)
	return err
}

func (self AppService) FindHotKeywords() (ps []entity.AppParams) {
	_ = self.db.Where("type=?", entity.ATHk).Desc("updated_at").Find(&ps)
	return
}

func (self AppService) DeleteHotKeyword(id int64) (err error) {
	var ad entity.AppParams
	_, err = self.db.Id(id).Get(&ad)
	if err != nil {
		return
	}
	_, err = self.db.Delete(&ad)
	return err
}

func (self AppService) SetFirstHotKeyword(id int64) (err error) {
	var ad entity.AppParams
	_, err = self.db.Id(id).Get(&ad)
	if err != nil {
		return
	}
	_, err = self.db.Id(ad.Id).Update(&ad)
	return err
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

func (self AppService) SaveSlogan(s entity.AppParams) (err error) {
	if s.Id == 0 { //new
		s.Type = entity.ATSg
		_, err = self.db.Insert(&s)
		return
	} else {
		p, exists := self.GetSlogan()
		if exists {
			p.Value = s.Value
			_, err = self.db.Id(p.Id).Cols("value").Update(&p)
		}
		return
	}
}

func (self AppService) SavePayment(ps ...entity.Payment) (err error) {
	_, err = self.db.Insert(&ps)
	return
}

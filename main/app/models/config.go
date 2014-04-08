package models

import (
	. "github.com/itang/gotang"
	"github.com/lunny/xorm"
	//"github.com/revel/revel"
	"github.com/itang/yunshang/main/app/models/entity"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type AppConfigService struct {
	db *xorm.Session
}

func NewAppConfigService(db *xorm.Session) *AppConfigService {
	return &AppConfigService{db}
}

func (self AppConfigService) InitData() error {
	for _, o := range entity.DefaultAppConfs {
		self.SaveOrUpdateConfigObject(o)
	}
	return nil
}

func (self AppConfigService) GetConfig(key string) (c entity.AppConfig, exists bool) {
	Assert(len(key) > 0, "Key 不能为空")

	exists, err := self.db.Where("key=?", key).Get(&c)
	AssertNoError(err, "GetConfig")

	return
}

func (self AppConfigService) GetOrSaveConfig(key string, value string) (c entity.AppConfig) {
	Assert(len(key) > 0, "Key 不能为空")

	c, exists := self.GetConfig(key)
	if !exists {
		return self.SaveOrUpdateConfig(key, value, "")
	}

	return
}

func (self AppConfigService) SaveOrUpdateConfig(key string, value string, desc string) (c entity.AppConfig) {
	return self.SaveOrUpdateConfigObject(entity.AppConfig{Key: key, Value: value, Description: desc})
}

func (self AppConfigService) SaveOrUpdateConfigObject(o entity.AppConfig) (c entity.AppConfig) {
	Assert(len(o.Key) > 0, "Key 不能为空")

	c, exists := self.GetConfig(o.Key)
	if !exists { //不存在
		c = o
		_, err := self.db.Insert(&c)
		AssertNoError(err, "SaveConfig")
	} else { //已经存在， 执行更新逻辑
		c.Value = o.Value
		session := self.db.Cols("value")
		if c.Description != o.Description && o.Description != "" {
			c.Description = o.Description
			session = self.db.Cols("value", "description")
		}
		_, err := session.Id(c.Id).Update(&c)
		AssertNoError(err, "SaveConfig")
	}

	return
}

func (self AppConfigService) FindConfigsBySection(section string) (ps []entity.AppConfig) {
	_ = self.db.Where("key like ?", section+"%").Asc("id").Find(&ps)
	return
}

func (self AppConfigService) FindAllConfigs() (ps []entity.AppConfig) {
	_ = self.db.Asc("id").Find(&ps)
	return
}

func (self AppConfigService) FindAllConfigsAsMap() map[string]entity.AppConfig {
	ret := make(map[string]entity.AppConfig)
	cs := self.FindAllConfigs()
	for _, c := range cs {
		ret[c.Key] = c
	}
	return ret
}

package models

import (
	"log"

	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models/data"
	"github.com/itang/yunshang/main/app/models/entity"
	db_module "github.com/itang/yunshang/modules/db"
	"github.com/itang/yunshang/modules/oauth"
)

func init() {
	app.OnAppInit(initDb)
}

// 初始化数据库相关
func initDb() {
	log.Println("Sync tables")
	err1 := db_module.Engine.Sync(
		&entity.User{},
		&entity.UserLevel{}, &entity.UserWorkKind{}, &entity.Location{}, &entity.UserDetail{},
		&entity.CompanyType{}, &entity.CompanyMainBiz{}, &entity.CompanyDetailBiz{},
		&entity.Company{},
		&oauth.UserSocial{},
		&entity.LoginLog{},
	)
	if err1 != nil {
		log.Fatalf("%v\n", err1)
	}

	log.Println("Init data")
	// init data
	data.TryInitData(db_module.Engine)
}

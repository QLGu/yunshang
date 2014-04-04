package data

import (
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/data/migrates"
	_ "github.com/itang/yunshang/main/app/data/migrations"
)

func init() {
	app.OnAppInit(initDb)
}

// 初始化数据库相关
func initDb() {
	migrates.AppInit()
	migrates.DataIniter.DoMigrate()
}

package main

import (
	//_ "github.com/go-sql-driver/mysql"
	"github.com/itang/gotang"

	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
)

var (
	driver = "postgres"
	spec   = "dbname=yunshangdb user=dbuser password=dbuser sslmode=disable"

	//#db.driver=mysql
	//#db.spec="dbuser:dbuser@/yunshangdb?charset=utf8"
)

func main() {
	Engine, err := xorm.NewEngine(driver, spec)
	gotang.AssertNoError(err, Engine.Ping())
	defer Engine.Close()

	Engine.ShowSQL = true

	dropTables(Engine)
}

// 删除应用创建所有的表
func dropTables(engine *xorm.Engine) {
	tables := []string{
		"t_company", "t_company_detail_biz",
		"t_company_main_biz",
		"t_company_type",
		"t_location",
		"t_user",
		"t_user_detail",
		"t_user_level",
		"t_user_work_kind",
		"t_user_social",
		"t_login_log",
	}
	for _, t := range tables {
		engine.Exec("drop table " + t)
	}
}

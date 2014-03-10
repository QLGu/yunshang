package main

import (
	"fmt"

	"github.com/itang/gotang"
	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
)

var (
	driver = "postgres"
	spec   = "dbname=yunshangdb user=dbuser password=dbuser sslmode=disable"
)

func main() {
	Engine, err := xorm.NewEngine(driver, spec)
	gotang.AssertNoError(err, "")
	defer Engine.Close()

	dropTables(Engine)
}

// 删除应用创建所有的表
func dropTables(engine *xorm.Engine) {
	tables := []string{
		"t_user_detail",
		"t_user_level",
		"t_user_work_kind",
		"t_user_social",
		"t_user",
		"t_login_log",
		"t_job_log",
		"t_delivery_address",
		"t_product",
		"t_product_price_rule",
	}
	for _, t := range tables {
		sql := fmt.Sprintf("drop table IF EXISTS %s CASCADE", t)
		_, err := engine.Exec(sql)
		fmt.Printf("%s, err: %v\n", sql, err)
	}
	fmt.Println("")
}

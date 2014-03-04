package db

import (
	"database/sql"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itang/gotang"
	"github.com/itang/reveltang"
	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

var (
	Db     *sql.DB
	Engine *xorm.Engine
)

func ModuleInit() {
	log.Printf("Init Module %v", "db")

	driver, spec := driverInfoFromConfig()
	revel.INFO.Printf("driver: %s, spec: %s", driver, hidePasswordFromSpec(spec))

	db, err := sql.Open(driver, spec)
	gotang.AssertNoError(err)

	engine, err := xorm.NewEngine(driver, spec)
	gotang.AssertNoError(err, engine.Ping())

	engine.SetTableMapper(xorm.NewPrefixMapper(xorm.SnakeMapper{}, "t_"))
	engine.ShowSQL = revel.Config.BoolDefault("db.show_sql", false)

	Db, Engine = db, engine
}

func driverInfoFromConfig() (driver string, spec string) {
	driver = reveltang.ForceGetConfig("db.driver")
	spec = reveltang.ForceGetConfig("db.spec")
	return
}

func hidePasswordFromSpec(spec string) string {
	re1 := regexp.MustCompile(" password=(.*) ") // postgres
	re2 := regexp.MustCompile(":.*@")            // mysql
	return re2.ReplaceAllString(re1.ReplaceAllString(spec, " password=****** "), ":******@")
}

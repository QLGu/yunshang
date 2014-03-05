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

// 模块初始化入口
// 主要初始化Db *sql.DB和Engine *xorm.Engine
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

// 执行Db操作
// session由Engine自动New
func Do(call func(*xorm.Session) error) error {
	gotang.Assert(Engine != nil, "db.Engine must init before use!")

	return DoWithSession(Engine.NewSession(), call)
}

// 执行Db操作
// 参数session： xorm会话
// 参数call： 基于session的Db操作函数
func DoWithSession(session *xorm.Session, call func(*xorm.Session) error) error {
	session.Begin()
	defer func() {
		if err := recover(); err != nil {
			session.Rollback()
		}
		session.Close()
	}()
	if err := call(session); err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

////////////////////////////////////////////////////////////////////////////////////
// private
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

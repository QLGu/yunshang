package app

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itang/gotang"
	gotang_time "github.com/itang/gotang/time"
	"github.com/itang/yunshang/main/app/models/entity"
	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
)

var (
	Db     *sql.DB
	Engine *xorm.Engine
)

func init() {

	initRevelFilter()

	revel.OnAppStart(func() {
		initDb()
	})
}

/////////////////////////////////////////////////////
func initRevelFilter() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.TemplateFuncs["webTitle"] = func(prefix string) (webTitle string) {
		const KEY = "cache.web.title"
		if err := cache.Get(KEY, &webTitle); err != nil {
			webTitle = forceGetConfig("web.title")
			go cache.Set(KEY, webTitle, 24*30*time.Hour)
		}
		return
	}

	revel.TemplateFuncs["logined"] = func(session revel.Session) bool {
		v, e := session["user"]
		return e == true && v != ""
	}
}

func initDb() {
	// Db
	// Read configuration.
	driver, spec := driverInfoFromConfig()
	var err error
	Db, err = sql.Open(driver, spec)
	gotang.AssertNoError(err)
	//gotang.AssertNoError(Db.Ping())

	Engine, err = xorm.NewEngine(driver, spec)
	gotang.AssertNoError(err)
	gotang.AssertNoError(Engine.Ping())

	Engine.SetTableMapper(xorm.NewPrefixMapper(xorm.SnakeMapper{}, "t_"))
	Engine.ShowSQL = revel.Config.BoolDefault("db.show_sql", false)

	err1 := Engine.Sync(
		&entity.User{}, &entity.UserLevel{}, &entity.UserWorkKind{}, &entity.Location{}, &entity.UserDetail{},
		&entity.CompanyType{}, &entity.CompanyMainBiz{}, &entity.CompanyDetailBiz{},
		&entity.Company{})
	if err1 != nil {
		log.Fatalf("%v\n", err1)
	}

	// init data
	tryInitData()
}

func driverInfoFromConfig() (driver string, spec string) {
	driver = forceGetConfig("db.driver")
	spec = forceGetConfig("db.spec")
	return
}

func hidePassword(spec string) string {
	re1 := regexp.MustCompile(" password=(.*) ") // postgres
	re2 := regexp.MustCompile(":.*@")            // mysql
	return re2.ReplaceAllString(re1.ReplaceAllString(spec, " password=****** "), ":******@")
}

func tryInitData() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uid := fmt.Sprintf("itang-%s-%d", gotang_time.FormatDefault(time.Now()), r.Int())

	users := []entity.User{{LoginName: uid, RealName: uid}}

	for _, user := range users {
		_, err := Engine.Insert(&user)
		gotang.AssertNoError(err)
	}
}

func forceGetConfig(key string) (value string) {
	var exists bool
	value, exists = revel.Config.String(key)
	gotang.Assert(exists, key)
	log.Printf("%s: %v", key, value)
	return
}

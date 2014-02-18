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
	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
	"github.com/robfig/revel"

	"github.com/itang/yunshang/main/app/models/entity"
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

	err1 := Engine.Sync(&entity.User{})
	if err1 != nil {
		log.Fatalf("%v\n", err1)
	}

	// init data
	tryInitData()
}

func driverInfoFromConfig() (driver string, spec string) {
	var exists bool
	driver, exists = revel.Config.String("db.driver")
	gotang.Assert(exists, "db.driver")
	log.Printf("db.driver: %v", driver)

	spec, exists = revel.Config.String("db.spec")
	gotang.Assert(exists, "db.spec")
	log.Printf("db.spec: %v", hidePassword(spec))

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

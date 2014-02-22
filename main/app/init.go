package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dchest/captcha"
	_ "github.com/go-sql-driver/mysql"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	_ "github.com/lib/pq"
	"github.com/lunny/xorm"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
)

var (
	Db     *sql.DB
	Engine *xorm.Engine
)

const (
	// Standard width and height of a captcha image.
	CaptchaWidth  = 120
	CaptchaHeight = 40
)

func init() {

	initRevelFilter()

	initRevelTemplateFuncs()

	revel.OnAppStart(installHandlers)

	revel.OnAppStart(initDb)
}

/////////////////////////////////////////////////////
func initRevelFilter() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		revel.RouterFilter, // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter, // Parse parameters into Controller.Params.
		revel.SessionFilter, // Restore and write the session cookie.
		revel.FlashFilter, // Restore and write the flash cookie.
		revel.ValidationFilter, // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter, // Resolve the requested language
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter, // Compress the result.
		revel.ActionInvoker, // Invoke the action.
	}
}

func initRevelTemplateFuncs() {
	revel.TemplateFuncs["emptyOr"] = func(value interface{}, other interface{}) interface{} {
		switch value.(type) {
		case string:
		{
			s, _ := value.(string)
			if s == "" {
				return other
			}
		}
		}
		if value == nil {
			return other
		}
		return value
	}

	revel.TemplateFuncs["webTitle"] = func(prefix string) (webTitle string) {
		const KEY = "cache.web.title"
		if err := cache.Get(KEY, &webTitle); err != nil {
			webTitle = forceGetConfig("web.title")
			go cache.Set(KEY, webTitle, 24*30*time.Hour)
		}
		return
	}

	revel.TemplateFuncs["urlWithHost"] = func(value string) string {
		host := revel.Config.StringDefault("web.host", "localhost:9000")
		return "http://" + host + value
	}

	revel.TemplateFuncs["logined"] = func(session revel.Session) bool {
		_, ok := session["user"]
		return ok
	}

	revel.TemplateFuncs["isAdmin"] = func(session revel.Session) bool {
		user, _ := session["user"]
		// TODO
		return user == "admin"
	}

	revel.TemplateFuncs["isAdminByName"] = func(name string) bool {
		// TODO
		return name == "admin"
	}

	revel.TemplateFuncs["valueAsName"] = func(value interface{}, theType string) string {
		switch theType {
		case "user_enabled":
		{
			v := fmt.Sprintf("%v", value)
			if v == "true" {
				return "激活/有效"
			} else {
				return "未激活/禁用"
			}
		}

		default:
			return ""
		}
	}

	revel.TemplateFuncs["valueOppoAsName"] = func(value interface{}, theType string) string {
		switch theType {
		case "user_enabled":
		{
			v := fmt.Sprintf("%v", value)
			if v == "false" {
				return "激活"
			} else {
				return "禁用"
			}
		}

		default:
			return ""
		}
	}
}

func installHandlers() {
	var (
		serveMux     = http.NewServeMux()
		revelHandler = revel.Server.Handler
	)
	serveMux.Handle("/", revelHandler)
	serveMux.Handle("/captcha/", captcha.Server(CaptchaWidth, CaptchaHeight))
	revel.Server.Handler = serveMux
}

func initDb() {
	// Db
	// Read configuration.
	driver, spec := driverInfoFromConfig()
	var err error
	Db, err = sql.Open(driver, spec)
	gotang.AssertNoError(err)
	gotang.AssertNoError(Db.Ping())

	Engine, err = xorm.NewEngine(driver, spec)
	gotang.AssertNoError(err)
	//gotang.AssertNoError(Engine.Ping())

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
	total, _ := Engine.Where("login_name = ?", "admin").Count(&entity.User{})
	if total == 0 {
		revel.INFO.Printf("init data")
		admin := entity.User{Email: "livetang@qq.com",
			CryptedPassword: utils.Sha1("computer"), LoginName: "admin",
			Enabled: true}
		users := []entity.User{admin}
		for _, user := range users {
			_, err := Engine.Insert(&user)
			gotang.AssertNoError(err)
		}
	}
}

func forceGetConfig(key string) (value string) {
	var exists bool
	value, exists = revel.Config.String(key)
	gotang.Assert(exists, key)
	log.Printf("%s: %v", key, value)
	return
}

package app

import (
	"net/http"

	"github.com/dchest/captcha"
	db_module "github.com/itang/yunshang/modules/db"
	mail_module "github.com/itang/yunshang/modules/mail"
	"github.com/revel/revel"
)

const (
	// Standard width and height of a captcha image.
	captchaWidth  = 120
	captchaHeight = 40
)

var appStartupHooks []func()

func OnAppInit(f func()) {
	appStartupHooks = append(appStartupHooks, f)
}

func init() {
	initRevelFilter()

	revel.OnAppStart(installHandlers)

	revel.OnAppStart(initModules)

	revel.OnAppStart(func() {
		for _, hook := range appStartupHooks {
			hook()
		}
	})
}

/////////////////////////////////////////////////////
// 初始化Revel的过滤器
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

// 装载Http handlers
func installHandlers() {
	var (
		serveMux     = http.NewServeMux()
		revelHandler = revel.Server.Handler
	)
	serveMux.Handle("/", revelHandler)
	serveMux.Handle("/captcha/", captcha.Server(captchaWidth, captchaHeight))
	revel.Server.Handler = serveMux
}

// 初始化各个模块
func initModules() {
	db_module.ModuleInit()
	mail_module.ModuleInit()
}

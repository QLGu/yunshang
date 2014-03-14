package controllers

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/itang/gotang"
	"github.com/itang/reveltang"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/itang/yunshang/modules/oauth"
	"github.com/itang/yunshang/modules/oauth/apps"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

var SocialAuth *oauth.SocialAuth

func init() {
	revel.ERROR_CLASS = "error"

	revel.InterceptMethod((*XOrmController).begin, revel.BEFORE)

	revel.InterceptMethod((*XOrmTnController).begin, revel.BEFORE)
	revel.InterceptMethod((*XOrmTnController).commit, revel.AFTER)
	revel.InterceptMethod((*XOrmTnController).rollback, revel.PANIC)

	revel.InterceptMethod((*AppController).init, revel.BEFORE)

	revel.InterceptMethod((*ShouldLoginedController).checkUser, revel.BEFORE)
	revel.InterceptMethod((*AdminController).checkAdminUser, revel.BEFORE)

	initRevelTemplateFuncs()

	app.OnAppInit(initOAuth)
}

func initOAuth() {
	log.Println("Init OAuth")

	var clientId, secret string
	appURL := revel.Config.StringDefault("social_auth_url", "http://"+revel.Config.StringDefault("web.host", ""))
	if len(appURL) > 0 {
		oauth.DefaultAppUrl = appURL
	}

	clientId = revel.Config.StringDefault("weibo_client_id", "")
	secret = revel.Config.StringDefault("weibo_client_secret", "")
	gotang.Assert(clientId != "" && secret != "", "weibo_client_id和weibo_client_secret不能为空")

	err := oauth.RegisterProvider(apps.NewWeibo(clientId, secret))
	gotang.AssertNoError(err, "")

	//clientId = revel.Config.StringDefault("qq_client_id","")
	//secret = revel.Config.StringDefault("qq_client_secret","")
	//err = oauth.RegisterProvider(apps.NewQQ(clientId, secret))

	SocialAuth = oauth.NewSocial("/passport/open/", new(socialAuther))
}

// 初始化Revel 模板的Functions
func initRevelTemplateFuncs() {
	log.Println("Init Revel Template Functions")

	revel.TemplateFuncs["inc"] = func(i1, i2 int) int {
		return i1 + i2
	}

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
			webTitle = reveltang.ForceGetConfig("web.title")
			go cache.Set(KEY, webTitle, 24*30*time.Hour)
		}
		return
	}

	revel.TemplateFuncs["urlWithHost"] = func(value string) string {
		host := revel.Config.StringDefault("web.host", "localhost:9000")
		return "http://" + host + value
	}

	revel.TemplateFuncs["logined"] = func(session revel.Session) bool {
		_, ok := session["uid"]
		return ok
	}

	revel.TemplateFuncs["isAdmin"] = func(session revel.Session) bool {
		user, _ := session["screen_name"]
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

		case "user_gender":
			{
				v := fmt.Sprintf("%v", value)
				switch v {
				case "male":
					return "男"
				case "female":
					return "女"
				default:
					return ""
				}
			}
		case "company_type":
			{
				v := fmt.Sprintf("%v", value)
				switch v {
				case "1":
					return "企业单位"
				case "2":
					return "个体经营"
				case "3":
					return "事业单位或社会团体"
				default:
					return ""
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

	revel.TemplateFuncs["siteYear"] = func(_ string) string {
		sy := "2013"
		ny := time.Now().Format("2006")
		return sy + "-" + ny
	}

	revel.TemplateFuncs["active"] = func(s1, s2 string) string {
		if strings.HasPrefix(s2, s1) {
			return "active"
		}
		return ""
	}

	revel.TemplateFuncs["radiox"] = func(f *revel.Field, val string, rval string) template.HTML {
		checked := ""
		if f.Flash() == val {
			checked = " checked"
		} else if rval == val {
			checked = " checked"
		}
		return template.HTML(fmt.Sprintf(`<input type="radio" name="%s" value="%s"%s>`,
			html.EscapeString(f.Name), html.EscapeString(val), checked))
	}

	revel.TemplateFuncs["checkboxx"] = func(f *revel.Field, val string, rval string) template.HTML {
		checked := ""
		if f.Flash() == val {
			checked = " checked"
		} else if rval == val {
			checked = " checked"
		}
		return template.HTML(fmt.Sprintf(`<input type="checkbox" name="%s" value="%s"%s>`,
			html.EscapeString(f.Name), html.EscapeString(val), checked))
	}

	revel.TemplateFuncs["optionx"] = func(f *revel.Field, val, label string, rval string) template.HTML {
		selected := ""
		if f.Flash() == val {
			selected = " selected"
		} else if rval == val {
			selected = " selected"
		}
		return template.HTML(fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			html.EscapeString(val), selected, html.EscapeString(label)))
	}

	revel.TemplateFuncs["flash"] = func(renderArgs map[string]interface{}, name string) string {
		v, _ := renderArgs["flash"].(map[string]string)[name]
		return v
	}

	revel.TemplateFuncs["levelName"] = func(user entity.User) string {
		var ret string
		_ = db.Do(func(session *xorm.Session) (err error) {
			userLevel, ok := models.NewUserService(session).GetUserLevel(&user)
			if !ok {
				return fmt.Errorf("Get Nothing UserLevel")
			}
			ret = userLevel.Name
			return
		})
		return ret
	}

	revel.TemplateFuncs["categoryChildren"] = func(id int64) (ps []entity.ProductCategory) {
		db.Do(func(session *xorm.Session) (err error) {
			ps = models.NewProductService(session).FindAllAvailableCategoriesByParentId(id)
			return
		})
		return
	}

	revel.TemplateFuncs["zeroAsEmpty"] = func(v interface{}) interface{} {
		switch v.(type) {
		case int, int32, int64:
			if v == 0 {
				return ""
			}
		case time.Time:
			if v.(time.Time).IsZero() {
				return ""
			}
		}
		return v
	}
}

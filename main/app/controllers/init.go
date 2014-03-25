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
	"strconv"
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

	var ystTemplateFuncs = map[string]interface{}{
		"inc": func(i1, i2 int) int {
			return i1 + i2
		},
		"add": func(i1, i2 int) int {
			return i1 + i2
		},
		"sub": func(i1, i2 int) int {
			return i1 - i2
		},
		"emptyOr": func(value interface{}, other interface{}) interface{} {
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
		},
		"webTitle": func(prefix string) (webTitle string) {
			const KEY = "cache.web.title"
			if err := cache.Get(KEY, &webTitle); err != nil {
				webTitle = reveltang.ForceGetConfig("web.title")
				go cache.Set(KEY, webTitle, 24*30*time.Hour)
			}
			return
		},
		"urlWithHost": func(value string) string {
			host := revel.Config.StringDefault("web.host", "localhost:9000")
			return "http://" + host + value
		},
		"logined": func(session revel.Session) bool {
			_, ok := session["uid"]
			return ok
		},
		"isAdmin": func(session revel.Session) bool {
			user, _ := session["screen_name"]
			// TODO
			return user == "admin"
		},
		"isAdminByName": func(name string) bool {
			// TODO
			return name == "admin"
		},
		"valueAsName": func(value interface{}, theType string) string {
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
		},
		"valueOppoAsName": func(value interface{}, theType string) string {
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
		},
		"siteYear": func(_ string) string {
			sy := "2013"
			ny := time.Now().Format("2006")
			return sy + "-" + ny
		},
		"active": func(s1, s2 string) string {
			if strings.HasPrefix(s2, s1) {
				return "active"
			}
			return ""
		},
		"current": func(s1, s2 string) string {
			if strings.HasPrefix(s2, s1) {
				return "current"
			}
			return ""
		},
		"startsWith": strings.HasPrefix,
		"radiox": func(f *revel.Field, val string, rval string) template.HTML {
			checked := ""
			if f.Flash() == val {
				checked = " checked"
			} else if rval == val {
				checked = " checked"
			}
			return template.HTML(fmt.Sprintf(`<input type="radio" name="%s" value="%s"%s>`,
				html.EscapeString(f.Name), html.EscapeString(val), checked))
		},
		"checkboxx": func(f *revel.Field, val string, rval string) template.HTML {
			checked := ""
			if f.Flash() == val {
				checked = " checked"
			} else if rval == val {
				checked = " checked"
			}
			return template.HTML(fmt.Sprintf(`<input type="checkbox" name="%s" value="%s"%s>`,
				html.EscapeString(f.Name), html.EscapeString(val), checked))
		},
		"optionx": func(f *revel.Field, val, label string, rval string) template.HTML {
			selected := ""
			if f.Flash() == val {
				selected = " selected"
			} else if rval == val {
				selected = " selected"
			}
			return template.HTML(fmt.Sprintf(`<option value="%s"%s>%s</option>`,
				html.EscapeString(val), selected, html.EscapeString(label)))
		},
		"flash": func(renderArgs map[string]interface{}, name string) string {
			v, _ := renderArgs["flash"].(map[string]string)[name]
			return v
		},
		"levelName": func(user entity.User) string {
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
		},
		"zeroAsEmpty": func(v interface{}) interface{} {
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
		},
		"renderArgs": func(key string, renderArgs map[string]interface{}) interface{} {
			v, ok := renderArgs[key]
			if !ok {
				return ""
			}
			return v
		},
		"ys_top_categories": func(renderArgs map[string]interface{}) (ps []entity.ProductCategory) {
			db.DoWithSession(xormSession(renderArgs), func(session *xorm.Session) error {
				ps = models.NewProductService(session).FindAvailableTopCategories()
				return nil
			})
			return
		},
		"ys_category_children": func(id int64, renderArgs map[string]interface{}) (ps []entity.ProductCategory) {
			db.DoWithSession(xormSession(renderArgs), func(session *xorm.Session) error {
				ps = models.NewProductService(session).FindAllAvailableCategoriesByParentId(id)
				return nil
			})
			return
		},
		"ys_recommend_providers": func(renderArgs map[string]interface{}) (ps []entity.Provider) {
			db.DoWithSession(xormSession(renderArgs), func(session *xorm.Session) error {
				ps = models.NewProductService(session).RecommendProviders()
				return nil
			})
			return
		},
		"ys_carts": func(renderArgs map[string]interface{}) (ret int64) {
			uid, ok := uidFromSession(renderArgs)
			if !ok {
				return 0
			}

			db.DoWithSession(xormSession(renderArgs), func(session *xorm.Session) error {
				ret = models.NewOrderService(session).UserCarts(uid)
				return nil
			})
			return
		},
		"ys_can_buy": func(p entity.Product) bool {
			return p.Enabled && p.StockNumber > 0 && p.MinNumberOfOrders <= p.StockNumber
		},
		"boolStr": func(v bool) string {
			if v {
				return "true"
			}
			return "false"
		},
		"mod": func(i int, j int) int {
			return i % j
		},
		"rawjs": func(s string) template.JS {
			return template.JS(s)
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"newline": func(index int, maxline int) bool {
			i := index + 1
			return i%maxline == 1 && i != 1
		},
		"notEq": func(a, b interface{}) bool {
			return !revel.Equal(a, b)
		},
		"ys_slogan": func(renderArgs map[string]interface{}) (ret string) {
			db.DoWithSession(xormSession(renderArgs), func(session *xorm.Session) error {
				ret = models.NewAppService(session).GetSloganContent()
				return nil
			})
			return
		},
		"truncStr": func(s string, le int, a string) string {
			if ulen(s) < le {
				return substr(s, 0, le)
			}
			return substr(s, 0, le) + a
		},
	}

	for k, v := range ystTemplateFuncs {
		_, exists := revel.TemplateFuncs[k]
		gotang.Assert(!exists, "不能覆盖已有TemplateFuncs!")
		revel.TemplateFuncs[k] = v
	}
}

func xormSession(renderArgs map[string]interface{}) *xorm.Session {
	session, exists := renderArgs["_db"]
	gotang.Assert(exists, `renderArgs["_db"] 不存在`)
	return session.(*xorm.Session)

}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func ulen(s string) int {
	runes := []rune(s)
	return len(runes)
}

func uidFromSession(renderArgs map[string]interface{}) (int64, bool) {
	session := renderArgs["session"]
	s, ok := session.(revel.Session)
	gotang.Assert(ok, "")

	uid, ok := s["uid"]
	if !ok {
		return 0, false
	}

	id, err := strconv.Atoi(uid)
	if err != nil {
		return 0, false
	}
	return int64(id), true
}

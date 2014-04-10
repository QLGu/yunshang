package controllers

import (
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/itang/gotang"
	gtemplate "github.com/itang/gotang/template"
	grtemplate "github.com/itang/reveltang/template"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/cache"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/itang/yunshang/modules/oauth"
	"github.com/itang/yunshang/modules/oauth/apps"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
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
	appURL := cache.GetConfig("site.basic.host")
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
		"addint64": func(i1 int64, i2 int) int {
			return int(i1) + i2
		},
		"urlWithHost": func(value string) string {
			host := cache.GetConfig("site.basic.host")
			return "http://" + host + value
		},
		"logined": func(session revel.Session) bool {
			_, ok := session["uid"]
			return ok
		},
		"not_logined": func(session revel.Session) bool {
			_, ok := session["uid"]
			return !ok
		},
		"isManager":      isManager,
		"isAdmin":        isAdmin,
		"isSellManager":  isSellManager,
		"isSuperManager": isSuperAdmin,
		"hasRole":        hasRole,
		"isSuperRole": func(session revel.Session) bool {
			return hasRole("超级管理员", session)
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
		"siteYear": func() string {
			sy := revel.Config.StringDefault("site.basic.start_year", "2014")
			ny := time.Now().Format("2006")
			if sy == ny {
				return sy
			}
			return sy + "-" + ny
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
		"ys_ad_images":             cache.GetAdImages,
		"ys_latest_news":           cache.GetLatestNews,
		"ys_pref_products":         cache.GetPrefProducts,
		"ys_hot_keywords":          cache.GetHotKeywords,
		"ys_top_categories":        cache.GetTopCategories,
		"ys_category_children":     cache.GetCategoryChildren,
		"ys_recommend_providers":   cache.GetRecommendProviders,
		"ys_latest_products":       cache.GetLatestProducts,
		"ys_specialoffer_products": cache.GetSpecialofferProducts,
		"ys_hot_products":          cache.GetHotProducts,
		"ys_service_categories":    cache.GetServiceCategories,
		"ys_config":                cache.GetConfig,
		"ys_slogan":                cache.GetSloganContent,
		"ys_news_by_category":      cache.GetNewsByCategory,
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
			return p.Enabled && p.StockNumber > 0 && p.MinNumberOfOrders <= p.StockNumber && p.Price > 0
		},
		"ys_online_support_qq_as_json": cache.GetOnlineSupportQQAsJSON,
		"ys_has_filters_for_product": func(p int64, code string, q string) bool {
			return p != 0 || code != "" || q != ""
		},
		"ys_ie": func() template.HTML {
			return template.HTML(`<!--[if lt IE 9]>
  <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
  <script src="/public/media/js/excanvas.min.js"></script>
  <script src="/public/media/js/respond.min.js"></script>
  <![endif]-->`)
		},
		"ys_payment": func(id int64) (ret entity.Payment) {
			db.Do(func(session *xorm.Session) error {
				ret, _ = models.NewOrderService(session).GetPaymentById(id)
				return nil
			})
			return
		},
		"ys_login_providers": func() (ret []oauth.Provider) {
			return oauth.GetProviders()
		},
	}

	doMergeMap(revel.TemplateFuncs, ystTemplateFuncs, grtemplate.ExtTemplateFuncs, gtemplate.ExtTemplateFuncs)
}

func doMergeMap(target map[string]interface{}, froms ...map[string]interface{}) {
	for _, from := range froms {
		for k, v := range from {
			_, exists := target[k]
			gotang.Assert(!exists, "不能覆盖已有TemplateFuncs!")
			target[k] = v
		}
	}
}

func xormSession(renderArgs map[string]interface{}) *xorm.Session {
	session, exists := renderArgs["_db"]
	gotang.Assert(exists, `renderArgs["_db"] 不存在`)
	return session.(*xorm.Session)

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

func hasRole(role string, session revel.Session) bool {
	roles, ok := session["roles"]
	if !ok {
		return false
	}
	return strings.Contains(roles, "#"+role)
}

func isManager(session revel.Session) bool {
	return isAdmin(session) || isSellManager(session) || isSuperAdmin(session)
}

func isAdmin(session revel.Session) bool {
	return hasRole("管理员", session) || hasRole("超级管理员", session)
}

func isSuperAdmin(session revel.Session) bool {
	return hasRole("超级管理员", session)
}

func isSellManager(session revel.Session) bool {
	return hasRole("销售", session) || hasRole("超级管理员", session)
}

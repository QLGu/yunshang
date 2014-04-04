package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/itang/gotang"
	gtemplate "github.com/itang/gotang/template"
	grtemplate "github.com/itang/reveltang/template"
	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/cache"
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

func isAdmin(session revel.Session) bool {
	user, _ := session["screen_name"]
	// TODO
	return user == "admin"
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
		"urlWithHost": func(value string) string {
			host := revel.Config.StringDefault("web.host", "localhost:9000")
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
		"isAdmin": isAdmin,
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
		"ys_ad_images":cache.GetAdImages,
		"ys_latest_news":cache.GetLatestNews,
		"ys_pref_products":cache.GetPrefProducts,
		"ys_hot_keywords":cache.GetHotKeywords,
		"ys_top_categories":cache.GetTopCategories,
		"ys_category_children":cache.GetCategoryChildren,
		"ys_recommend_providers": cache.GetRecommendProviders,
		"ys_latest_products": cache.GetLatestProducts,
		"ys_specialoffer_products": cache.GetSpecialofferProducts,
		"ys_hot_products":cache.GetHotProducts,
		"ys_service_categories": cache.GetServiceCategories,
		"ys_config": cache.GetConfig,
		"ys_slogan":cache.GetSloganContent,
		"ys_news_by_category": cache.GetNewsByCategory,
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
	}

	doMergeMap(revel.TemplateFuncs, ystTemplateFuncs, grtemplate.ExtTemplateFuncs, gtemplate.ExtTemplateFuncs)
}

func doMergeMap(target map[string]interface{}, froms ... map[string]interface{}) {
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

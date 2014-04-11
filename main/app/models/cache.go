package models

import (
	"fmt"
	"strings"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

func init() {
	Emitter.On("update-cache", func(bean interface{}) {
		var t = utils.TypeOfTarget(bean).Name()
		cacheKeys := utils.GetCacheKeys()
		for _, k := range cacheKeys {
			if strings.HasSuffix(k, "_"+t) {
				revel.INFO.Printf("update cache %s of type %s", k, t)
				utils.ClearCache(k)
			}
		}

		//配置热更新
		if t == "AppConfig" {
			e, ok := bean.(*entity.AppConfig)
			if ok {
				switch {
				case strings.Contains(e.Key, "site.mail"):
					FireEvent(EventObject{Name: EReloadMailConfig})
				case strings.Contains(e.Key, "site.alipay"):
					FireEvent(EventObject{Name: EReloadAlipayConfig})
				}
			}
		}
	})
}

var CacheSystem = cacheSystem{}

type cacheSystem struct {
}

//应用配置
func (e cacheSystem) GetAppConfigs() (ret map[string]entity.AppConfig) {
	utils.Cache("ys_configs_AppConfig", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewAppConfigService(session).FindAllConfigsAsMap()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetConfig(key string) string {
	configs := e.GetAppConfigs()
	ac, exists := configs[key]
	gotang.Assert(exists, "配置不存在,"+key)

	return ac.Value
}

func (e cacheSystem) GetOnlineSupportQQAsJSON() string {
	var ret string = ""
	utils.Cache("ys_GetOnlineSupportQQAsJSON_AppConfig", &ret, func(key string) (data interface{}) {
		s := e.GetConfig("site.contact.online_support_qq")
		if len(s) == 0 {
			return "[]"
		}
		var qqArr = strings.Split(s, ",")
		var rets = make([]string, 0)
		for _, qq := range qqArr {
			nameValue := strings.Split(qq, ":")
			if len(nameValue) != 2 {
				continue
			}
			name := nameValue[0]
			value := nameValue[1]
			rets = append(rets, fmt.Sprintf(`{"name":"%s", "qq":"%s"}`, name, value))
		}
		return "[" + strings.Join(rets, ",") + "]"
	})
	return ret
}

//广告词
func (e cacheSystem) GetSloganContent() string {
	var ret string = ""
	utils.Cache("ys_slogan_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewAppService(session).GetSloganContent()
			return nil
		})
		return
	})
	return ret
}

func (e cacheSystem) GetServiceCategories() (ret []entity.NewsCategory) {
	utils.Cache("ys_GetServiceCategories_NewsCategory", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewNewsService(session).FindAllAvailableServiceCategories()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetHotProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetHotProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindHotProducts(limit)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetSpecialofferProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetSpecialofferProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindSpecialOfferProducts(limit)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetLatestProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetLatestProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindLatestProducts(limit)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetTopCategories() (ret []entity.ProductCategory) {
	utils.Cache("ys_GetTopCategories_ProductCategory", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindAvailableTopCategories()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetCategoryChildren(id int64) (ret []entity.ProductCategory) {
	utils.Cache(fmt.Sprintf("ys_GetCategoryChildren_%d_ProductCategory", id), &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindAllAvailableCategoriesByParentId(id)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetRecommendProviders() (ret []entity.Provider) {
	utils.Cache("ys_GetRecommendProviders_Provider", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).RecommendProviders()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetHotKeywords() (ret []entity.AppParams) {
	utils.Cache("ys_GetHotKeywords_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewAppService(session).FindHotKeywords()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetPrefProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetPrefProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewProductService(session).FindPrefProducts(limit)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetLatestNews(limit int) (ret []entity.News) {
	utils.Cache("ys_GetLatestNews_News", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewNewsService(session).FindNews("", limit)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetAdImages() (ret []entity.AppParams) {
	utils.Cache("ys_GetAdImages_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewAppService(session).FindAdImages()
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) GetNewsByCategory(ctId int64) (ret []entity.News) {
	utils.Cache(fmt.Sprintf("ys_GetNewsByCategory_%d_News", ctId), &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = NewNewsService(session).FindAllAvailableNewsByCategory(ctId)
			return nil
		})
		return
	})
	return
}

func (e cacheSystem) UrlWithHost(value string) string {
	host := e.GetConfig("site.basic.host")
	return "http://" + host + value
}

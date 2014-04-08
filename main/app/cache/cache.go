package cache

import (
	"fmt"
	"strings"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

func init() {
	models.Emitter.On("update-cache", func(t string) {
		revel.INFO.Println("update cache of type ", t)
		cacheKeys := utils.GetCacheKeys()
		for _, k := range cacheKeys {
			if strings.HasSuffix(k, "_"+t) {
				utils.ClearCache(k)
			}
		}
	})
}

//应用配置
func GetAppConfigs() (ret map[string]entity.AppConfig) {
	utils.Cache("ys_configs_AppConfig", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewAppConfigService(session).FindAllConfigsAsMap()
			return nil
		})
		return
	})
	return
}

func GetConfig(key string) string {
	configs := GetAppConfigs()
	ac, exists := configs[key]
	gotang.Assert(exists, "配置不存在,"+key)

	return ac.Value
}

func GetOnlineSupportQQAsJSON() string {
	var ret string = ""
	utils.Cache("ys_GetOnlineSupportQQAsJSON_AppConfig", &ret, func(key string) (data interface{}) {
		s := GetConfig("site.contact.online_support_qq")
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
func GetSloganContent() string {
	var ret string = ""
	utils.Cache("ys_slogan_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewAppService(session).GetSloganContent()
			return nil
		})
		return
	})
	return ret
}

func GetServiceCategories() (ret []entity.NewsCategory) {
	utils.Cache("ys_GetServiceCategories_NewsCategory", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewNewsService(session).FindAllAvailableServiceCategories()
			return nil
		})
		return
	})
	return
}

func GetHotProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetHotProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindHotProducts(limit)
			return nil
		})
		return
	})
	return
}

func GetSpecialofferProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetSpecialofferProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindSpecialOfferProducts(limit)
			return nil
		})
		return
	})
	return
}

func GetLatestProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetLatestProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindLatestProducts(limit)
			return nil
		})
		return
	})
	return
}

func GetTopCategories() (ret []entity.ProductCategory) {
	utils.Cache("ys_GetTopCategories_ProductCategory", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindAvailableTopCategories()
			return nil
		})
		return
	})
	return
}

func GetCategoryChildren(id int64) (ret []entity.ProductCategory) {
	utils.Cache(fmt.Sprintf("ys_GetCategoryChildren_%d_ProductCategory", id), &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindAllAvailableCategoriesByParentId(id)
			return nil
		})
		return
	})
	return
}

func GetRecommendProviders() (ret []entity.Provider) {
	utils.Cache("ys_GetRecommendProviders_Provider", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).RecommendProviders()
			return nil
		})
		return
	})
	return
}

func GetHotKeywords() (ret []entity.AppParams) {
	utils.Cache("ys_GetHotKeywords_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewAppService(session).FindHotKeywords()
			return nil
		})
		return
	})
	return
}

func GetPrefProducts(limit int) (ret []entity.Product) {
	utils.Cache("ys_GetPrefProducts_Product", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewProductService(session).FindPrefProducts(limit)
			return nil
		})
		return
	})
	return
}

func GetLatestNews(limit int) (ret []entity.News) {
	utils.Cache("ys_GetLatestNews_News", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewNewsService(session).FindNews("", limit)
			return nil
		})
		return
	})
	return
}

func GetAdImages() (ret []entity.AppParams) {
	utils.Cache("ys_GetAdImages_AppParams", &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewAppService(session).FindAdImages()
			return nil
		})
		return
	})
	return
}

func GetNewsByCategory(ctId int64) (ret []entity.News) {
	utils.Cache(fmt.Sprintf("ys_GetNewsByCategory_%d_News", ctId), &ret, func(key string) (data interface{}) {
		db.Do(func(session *xorm.Session) error {
			data = models.NewNewsService(session).FindAllAvailableNewsByCategory(ctId)
			return nil
		})
		return
	})
	return
}

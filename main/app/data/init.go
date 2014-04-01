package data

import (
	"log"

	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models/entity"
	db_module "github.com/itang/yunshang/modules/db"
	"github.com/itang/yunshang/modules/oauth"
)

func init() {
	app.OnAppInit(initDb)
}

// 初始化数据库相关
func initDb() {
	log.Println("Sync tables")
	err1 := db_module.Engine.Sync(
		&entity.AppParams{},
		&entity.User{},
		&entity.UserLevel{},
		&entity.UserWorkKind{},
		&entity.UserDetail{},
		&oauth.UserSocial{},
		&entity.LoginLog{},
		&entity.JobLog{},
		&entity.DeliveryAddress{},
		&entity.ProductCategory{},
		&entity.Product{},
		&entity.ProductPrices{},
		&entity.ProductParams{},
		&entity.ProductStockLog{},
		&entity.ProductCollect{},
		&entity.Provider{},
		&entity.Cart{},
		&entity.Payment{},
		&entity.Order{},
		&entity.OrderDetail{},
		&entity.Shipping{},
		&entity.OrderLog{},
		&entity.Invoice{},
		&entity.Comment{},
		&entity.Inquiry{},
		&entity.InquiryReply{},
		&entity.NewsCategory{},
		&entity.News{},
		&entity.NewsParam{},
	)
	if err1 != nil {
		log.Fatalf("%v\n", err1)
	}

	log.Println("Init data")
	// init data
	go TryInitData(db_module.Engine)
}

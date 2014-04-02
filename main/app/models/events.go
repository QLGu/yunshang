package models

import (
	"github.com/chuckpreslar/emission"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/kr/pretty"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

var _emitter = emission.NewEmitter()

const (
	EStockLog       = "stock-log"
	EOrderLog       = "order-log"
	ECommon         = "common"
	EProductComment = "product-common"
	EPay            = "pay"
)

type EventObject struct {
	Name     string
	User     string
	UserId   int64
	SourceId int64
	Title    string
	Message  string
	Data     interface{}
}

func init() {
	_emitter.On(ECommon, func(e EventObject) {
		revel.INFO.Printf("%# v", pretty.Formatter(e))
		switch e.Name {
		case EStockLog:
			log := entity.ProductStockLog{ProductId: e.SourceId, User: e.User, Message: e.Message}
			db.Do(func(session *xorm.Session) error {
				_, err := session.Insert(&log)
				return err
			})

		case EOrderLog:
			log := entity.OrderLog{OrderId: e.SourceId, Message: e.Title + ": " + e.Message}
			db.Do(func(session *xorm.Session) error {
				_, err := session.Insert(&log)
				return err
			})

		case EProductComment:
			userId := e.UserId
			db.Do(func(session *xorm.Session) error {
				NewUserService(session).DoIncUserScores(userId, 1) // 评论每次1分
				return nil
			})

		case EPay:
			userId := e.UserId
			famount, ok := e.Data.(float64)
			gotang.Assert(ok, "data should be float64")
			amount := int(famount)
			db.Do(func(session *xorm.Session) error {
				NewUserService(session).DoIncUserScores(userId, amount/2) //2元 ， 1分
				return nil
			})

		default:
			revel.WARN.Println("Unknow Event", e)
		}
	})
}

func FireEvent(e EventObject) {
	_emitter.Emit(ECommon, e)
}

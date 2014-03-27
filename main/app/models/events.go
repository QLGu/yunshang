package models

import (
	"github.com/chuckpreslar/emission"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

var _emitter = emission.NewEmitter()

const (
	EStockLog = "stock-log"
	EOrderLog = "order-log"
	ECommon   = "common"
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
		default:
			revel.WARN.Println("Unknow Event", e)
		}
	})
}

func FireEvent(e EventObject) {
	_emitter.Emit(ECommon, e)
}

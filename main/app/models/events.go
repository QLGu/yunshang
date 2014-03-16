package models

import (
	"github.com/chuckpreslar/emission"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
)

var Emitter = emission.NewEmitter()

const (
	EStockLog = "stock-log"
)

func init() {
	Emitter.On(EStockLog, func(user string, productId int64, message string) {
		log := entity.ProductStockLog{ProductId: productId, User: user, Message: message}
		db.Do(func(session *xorm.Session) error {
			_, err := session.Insert(&log)
			return err
		})
	})
}

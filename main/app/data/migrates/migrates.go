package migrates

import (
	"log"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
)

var DataIniter = NewDataIniter()

func AppInit() {
	err := db.Engine.Sync(&entity.Migration{})
	gotang.AssertNoError(err, "Syunc Migration")
}

type Migration struct {
	Name string
	Desc string
	Do   func(session *xorm.Session) error
}

func NewDataIniter() *dataIniter {
	return &dataIniter{make(map[string]Migration)}
}

type dataIniter struct {
	ms map[string]Migration
}

func (self *dataIniter) RegistMigration(m Migration) {
	self.ms[m.Name] = m
}

func (self *dataIniter) DoMigrate() {
	db.Do(func(session *xorm.Session) error {
		appApi := models.NewAppService(session)
		existsMs := appApi.FindAllMigrationsAsMap()

		for name, m := range self.ms {
			_, exists := existsMs[name]
			if exists {
				log.Printf("Migration %s 已经存在， 跳过!", name)
				continue
			}

			log.Printf("运行: Migration %s ", name)
			err := (m.Do)(session)
			gotang.AssertNoError(err, name+" do, Has error!")

			err = appApi.SaveMigration(m.Name, m.Desc)
			gotang.AssertNoError(err, name+" do, Has error!")

		}
		return nil
	})
}

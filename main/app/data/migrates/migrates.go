package migrates

import (
	"fmt"
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
	gotang.AssertNoError(err, "Sync Migration")
}

type Migration struct {
	Name string
	Desc string
	Do   func(session *xorm.Session) error
}

func NewDataIniter() *dataIniter {
	return &dataIniter{make([]Migration, 0)}
}

type dataIniter struct {
	ms []Migration
}

func (self *dataIniter) RegistMigration(m Migration) {
	self.ms = append(self.ms, m)
}

func (self *dataIniter) DoMigrate() {
	for _, m := range self.ms {
		fmt.Println("migrate:", m.Name)
	}
	db.Do(func(session *xorm.Session) error {
		appApi := models.NewAppService(session)
		existsMs := appApi.FindAllMigrationsAsMap()

		for _, m := range self.ms {
			_, exists := existsMs[m.Name]
			if exists {
				log.Printf("Migration %s 已经存在， 跳过!", m.Name)
				continue
			}

			log.Printf("运行: Migration %s ", m.Name)
			err := (m.Do)(session)
			gotang.AssertNoError(err, m.Name+" do, Has error!")

			err = appApi.SaveMigration(m.Name, m.Desc)
			gotang.AssertNoError(err, m.Name+" do, Has error!")

			log.Printf("完成: Migration %s ", m.Name)
		}
		return nil
	})
}

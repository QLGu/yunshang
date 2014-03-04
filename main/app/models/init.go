package models

import (
	"log"
	"time"

	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models/data"
	"github.com/itang/yunshang/main/app/models/entity"
	db_module "github.com/itang/yunshang/modules/db"
	"github.com/itang/yunshang/modules/oauth"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"
)

func init() {
	app.OnAppInit(initDb)

	app.OnAppInit(func() {
		revel.INFO.Println("Jobs init")
		jobs.Schedule("0 19 21 * * *", jobs.Func(computeUsersScoresJob))
	})
}

// 初始化数据库相关
func initDb() {
	log.Println("Sync tables")
	err1 := db_module.Engine.Sync(
		&entity.User{},
		&entity.UserLevel{}, &entity.UserWorkKind{}, &entity.Location{}, &entity.UserDetail{},
		&entity.CompanyType{}, &entity.CompanyMainBiz{}, &entity.CompanyDetailBiz{},
		&entity.Company{},
		&oauth.UserSocial{},
		&entity.LoginLog{},
		&entity.JobLog{},
	)
	if err1 != nil {
		log.Fatalf("%v\n", err1)
	}

	log.Println("Init data")
	// init data
	data.TryInitData(db_module.Engine)
}

func computeUsersScoresJob() {
	jobName := "computeUsersScoresJob"

	revel.INFO.Printf("Job %s start... ", jobName)
	session := db_module.Engine.NewSession()

	defer func() {
		if err := recover(); err != nil {
			session.Rollback()
		}
	}()
	session.Begin()

	date := time.Now().Format("2006-01-02")

	js := NewJobService(session)
	if !js.ExistJobLog(jobName, date) {
		if err := DefaultUserService(session).ComputeUsersScores(date); err != nil {
			panic(err)
		}

		js.AddJobLog(jobName, date)
	}

	session.Commit()
	revel.INFO.Printf("Job %s end ", jobName)
}

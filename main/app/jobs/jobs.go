package jobs

import (
	"time"

	gtime "github.com/itang/gotang/time"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

type ComputeUsersScoresJob struct {
}

func (e ComputeUsersScoresJob) Run() {
	jobName := "computeUsersScoresJob"
	revel.INFO.Printf("Job %s start... ", jobName)

	db.Do(func(session *xorm.Session) error {
		date := gtime.RichTime(time.Now()).Yesterday().Format("2006-01-02")
		revel.INFO.Printf("computeUsersScoresJob date: %s", date)

		js := models.NewJobLogService(session)
		if !js.ExistJobLog(jobName, date) {
			if err := models.DefaultUserService(session).ComputeUsersScores(date); err != nil {
				panic(err)
			}

			js.AddJobLog(jobName, date)
		}
		return nil
	})

	revel.INFO.Printf("Job %s end ", jobName)
}

package jobs

import (
	"log"

	"github.com/itang/yunshang/main/app"
	"github.com/revel/revel/modules/jobs/app/jobs"
)

func init() {
	app.OnAppInit(func() {
		log.Printf("Init %s", "jobs")
		jobs.Schedule("cron.compute_scores", ComputeUsersScoresJob{})
	})
}

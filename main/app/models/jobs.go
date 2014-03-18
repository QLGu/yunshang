package models

import (
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

func NewJobLogService(db *xorm.Session) *JobLogService {
	return &JobLogService{db}
}

/////////////////////////////////////////////////

type JobLogService struct {
	db *xorm.Session
}

func (self JobLogService) ExistJobLog(name string, date string) bool {
	c, _ := self.db.Where("name=? and date=?", name, date).Count(&entity.JobLog{})
	return c > 0
}

func (self JobLogService) AddJobLog(name string, date string) (err error) {
	jobLog := entity.JobLog{Name: name, Date: date}
	_, err = self.db.Insert(&jobLog)
	return
}

package models

import (
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

type JobService interface {
	ExistJobLog(name string, date string) bool

	AddJobLog(name string, date string) error
}

func NewJobService(session *xorm.Session) JobService {
	return jobServiceImpl{session}
}

/////////////////////////////////////////////////

type jobServiceImpl struct {
	session *xorm.Session
}

func (self jobServiceImpl) ExistJobLog(name string, date string) bool {
	c, _ := self.session.Where("name=? and date=?", name, date).Count(&entity.JobLog{})
	return c > 0
}

func (self jobServiceImpl) AddJobLog(name string, date string) (err error) {
	jobLog := entity.JobLog{Name: name, Date: date}
	_, err = self.session.Insert(&jobLog)
	return
}

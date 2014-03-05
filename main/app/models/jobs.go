package models

import (
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

type JobLogService interface {
	ExistJobLog(name string, date string) bool

	AddJobLog(name string, date string) error
}

func NewJobLogService(session *xorm.Session) JobLogService {
	return jobLogServiceImpl{session}
}

/////////////////////////////////////////////////

type jobLogServiceImpl struct {
	session *xorm.Session
}

func (self jobLogServiceImpl) ExistJobLog(name string, date string) bool {
	c, _ := self.session.Where("name=? and date=?", name, date).Count(&entity.JobLog{})
	return c > 0
}

func (self jobLogServiceImpl) AddJobLog(name string, date string) (err error) {
	jobLog := entity.JobLog{Name: name, Date: date}
	_, err = self.session.Insert(&jobLog)
	return
}

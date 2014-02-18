package models

import (
	"github.com/itang/gotang"
	"github.com/lunny/xorm"

	"github.com/itang/yunshang/main/app/models/entity"
)

/*type User interface {
	GetId() int64
	GetLoginname() string
	GetRealname() string
}*/

type UserService interface {
	Total() int64

	FindAllUsers() []entity.User
}

func DefaultUserService(session *xorm.Session) UserService {
	return &defaultUserService{session}
}

type defaultUserService struct {
	session *xorm.Session
}

func (self defaultUserService) Total() int64 {
	ret, err := self.session.Count(&entity.User{})
	gotang.AssertNoError(err)

	return ret
}

func (self defaultUserService) FindAllUsers() (users []entity.User) {
	self.session.Find(&users)
	return
}

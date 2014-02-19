package models

import (
	"errors"
	"time"

	"github.com/itang/gotang"
	"github.com/lunny/xorm"

	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
)

type UserService interface {
	Total() int64

	FindAllUsers() []entity.User

	RegistUser(email, password string) (entity.User, error)

	ExistsUserByEmail(email string) bool

	Activate(email, code string) (entity.User, error)
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

func (self defaultUserService) RegistUser(email, password string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = password
	user.ActivationCode = utils.Uuid()
	user.ActivationCodeCreatedAt = time.Now()

	_, err = self.session.Insert(&user)
	return
}

func (self defaultUserService) ExistsUserByEmail(email string) bool {
	total, _ := self.session.Where("email=?", email).Count(&entity.User{})
	return total > 0
}

func (self defaultUserService) Activate(email, code string) (user entity.User, err error) {
	var users []entity.User
	err = self.session.Where("email=? and activation_code=?", email, code).Find(&users)
	if err != nil {
		return
	}
	if len(users) > 0 {
		user = users[0]
		user.Enabled = true
		user.ActivationCode = ""
		self.session.Id(user.Id).Cols("enabled", "activation_code").Update(&user)
		return
	} else {
		err = errors.New("no user find")
		return
	}
}

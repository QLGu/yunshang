package models

import (
	"errors"
	"time"

	"github.com/itang/gotang"
	"github.com/lunny/xorm"

	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
)

var UserTypeInstance = &entity.User{}

type SessionUser struct {
	Email     string
	LoginName string
	RealName  string
}

func (self SessionUser) DisplayName() string {
	if len(self.LoginName) == 0 {
		return self.Email
	}
	return self.LoginName
}

func ToSessionUser(user entity.User) SessionUser {
	return SessionUser{Email: user.Email, LoginName: user.LoginName, RealName: user.RealName}
}

type UserService interface {
	Total() int64

	FindAllUsers() []entity.User

	RegistUser(email, password string) (entity.User, error)

	ExistsUserByEmail(email string) bool

	Activate(email, code string) (entity.User, error)

	CheckUser(login, password string) (entity.User, bool)

	CheckUserByEmail(email string) (entity.User, bool)

	DoUserLogin(user *entity.User) error

	DoForgotPasswordApply(user *entity.User) error

	ResetUserPassword(email, code string) (string, error)
}

func DefaultUserService(session *xorm.Session) UserService {
	return &defaultUserService{session}
}

type defaultUserService struct {
	session *xorm.Session
}

func (self defaultUserService) Total() int64 {
	ret, err := self.session.Count(UserTypeInstance)
	gotang.AssertNoError(err)

	return ret
}

func (self defaultUserService) FindAllUsers() (users []entity.User) {
	self.session.Find(&users)
	return
}

func (self defaultUserService) RegistUser(email, password string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = utils.Sha1(password)
	user.ActivationCode = utils.Uuid()
	user.ActivationCodeCreatedAt = time.Now()

	_, err = self.session.Insert(&user)
	return
}

func (self defaultUserService) ExistsUserByEmail(email string) bool {
	total, _ := self.session.Where("email=?", email).Count(UserTypeInstance)
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
		err = errors.New("激活码不存在或已经失效！")
		return
	}
}

func (self defaultUserService) CheckUser(login, password string) (user entity.User, ok bool) {
	ok, err := self.session.Where("(email=? or login_name=?) and crypted_password = ?", login, login, utils.Sha1(password)).Get(&user)
	return user, ok && err == nil
}

func (self defaultUserService) CheckUserByEmail(email string) (user entity.User, ok bool) {
	ok, err := self.session.Where("email=?", email).Get(&user)
	return user, ok && err == nil
}

func (self defaultUserService) DoUserLogin(user *entity.User) error {
	user.LastSignAt = time.Now()
	_, err := self.session.Id(user.Id).Update(user)

	return err
}

func (self defaultUserService) DoForgotPasswordApply(user *entity.User) error {
	user.PasswordResetCode = utils.Uuid()
	_, err := self.session.Id(user.Id).Update(user)
	return err
}

func (self defaultUserService) ResetUserPassword(email, code string) (newPassword string, err error) {
	newPassword = utils.RandomString(6)
	var user entity.User
	ok, err := self.session.Where("email=? and password_reset_code = ?", email, code).Get(&user)
	if !ok {
		return "", errors.New("密码重置请求无效")
	}

	if err != nil {
		return "", err
	}

	user.CryptedPassword = utils.Sha1(newPassword)
	_, err = self.session.Id(user.Id).Cols("crypted_password").Update(&user)

	if err != nil {
		return "", err
	}
	return newPassword, err
}

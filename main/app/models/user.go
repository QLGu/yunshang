package models

import (
	"errors"
	"log"
	"time"

	"github.com/itang/gotang"
	gtime "github.com/itang/gotang/time"
	"github.com/lunny/xorm"

	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
)

var UserTypeInstance = &entity.User{}

var LoginLogTypeInstance = &entity.LoginLog{}

type SessionUser struct {
	Id        int64
	Email     string
	LoginName string
	RealName  string
	From      string
}

func (self SessionUser) DisplayName() string {
	if len(self.LoginName) == 0 {
		return self.Email
	}
	return self.LoginName
}

func ToSessionUser(user entity.User) SessionUser {
	from := user.From
	if len(from) == 0 {
		from = "Local"
	}
	return SessionUser{Id: user.Id, Email: user.Email, LoginName: user.LoginName, RealName: user.RealName, From: from}
}

type UserService interface {
	Total() int64

	GetUserById(id int64) (user entity.User, ok bool)

	FindAllUsers() []entity.User

	FindAllUsersForPage(ps *PageSearcher) PageData

	RegistUser(email, password string) (entity.User, error)

	ExistsUserByEmail(email string) bool

	Activate(email, code string) (entity.User, error)

	CheckUser(login, password string) (entity.User, bool)

	CheckUserByEmail(email string) (entity.User, bool)

	DoUserLogin(user *entity.User) error

	DoForgotPasswordApply(user *entity.User) error

	ResetUserPassword(email, code string) (string, error)

	GetUserByLogin(login string) (entity.User, bool)

	DoChangePassword(user *entity.User, rawPassword string) error

	VerifyPassword(cryptedPassword, rawPassword string) bool

	IsAdminUser(user *entity.User) bool

	ToggleUserEnabled(user *entity.User) error

	ConnectUser(id string, providerName string, email string) (entity.User, error)

	GetUserLevel(user *entity.User) (entity.UserLevel, bool)

	FindUserLevels() []entity.UserLevel

	FindUserLoginLogs(userId int64) []entity.LoginLog
}

func DefaultUserService(session *xorm.Session) UserService {
	return &defaultUserService{session}
}

//////////////////////////////////////////////////////////////////////////////
// impls
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
	user.Code = utils.Uuid()
	user.From = "Local"
	user.ActivationCodeCreatedAt = time.Now()

	_, err = self.session.Insert(&user)
	return
}

func (self defaultUserService) ConnectUser(id string, providerName string, email string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = ""
	user.ActivationCode = ""
	user.LoginName = providerName+id
	user.From = providerName
	user.Code = utils.Uuid()
	user.Enabled = true

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
	ok, err := self.session.Where("(email=? or login_name=?) and crypted_password = ? and enabled=?", login, login, utils.Sha1(password), true).Get(&user)
	return user, ok && err == nil
}

func (self defaultUserService) CheckUserByEmail(email string) (user entity.User, ok bool) {
	ok, err := self.session.Where("email=?", email).Get(&user)
	return user, ok && err == nil
}

// 用户成功登录之后
func (self defaultUserService) DoUserLogin(user *entity.User) error {
	// 更新最近一次登录时间
	user.LastSignAt = time.Now()
	_, err := self.session.Id(user.Id).Cols("last_sign_at").Update(user)

	date := user.LastSignAt.Format(gtime.ChinaDefaultDate)
	self.doUpdateLoginLog2(user.Id, date, user.LastSignAt)

	return err
}

func (self defaultUserService) doUpdateLoginLog(userId int64, date string, detailTime time.Time) (llog entity.LoginLog, new bool) {
	exists, err := self.session.Where("date = ?", date).And("user_id = ?", userId).Get(&llog)
	gotang.AssertNoError(err)

	if exists {
		llog.DetailTime = detailTime
		self.session.Id(llog.Id).Cols("detail_time").Update(&llog)
		new = false
	}else {
		llog.Date = date
		llog.DetailTime = detailTime
		llog.UserId = userId
		self.session.Insert(&llog)
		new = true
	}
	return
}

func (self defaultUserService) doUpdateLoginLog2(userId int64, date string, detailTime time.Time) (llog entity.LoginLog, new bool) {
	llog.Date = date
	llog.DetailTime = detailTime
	llog.UserId = userId
	self.session.Insert(&llog)
	new = true

	return
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

func (self defaultUserService) GetUserByLogin(login string) (user entity.User, ok bool) {
	ok, err := self.session.Where("email=? or login_name=?", login, login).Get(&user)
	return user, ok && err == nil
}

func (self defaultUserService) DoChangePassword(user *entity.User, rawPassword string) error {
	user.CryptedPassword = utils.Sha1(rawPassword)
	_, err := self.session.Id(user.Id).Cols("crypted_password").Update(user)
	return err
}

func (self defaultUserService) VerifyPassword(cryptedPassword, rawPassword string) bool {
	return cryptedPassword == utils.Sha1(rawPassword)
}

func (self defaultUserService) GetUserById(id int64) (user entity.User, ok bool) {
	ok, err := self.session.Id(id).Get(&user)
	return user, ok && err == nil
}

func (self defaultUserService) IsAdminUser(user *entity.User) bool {
	//TODO 改进判断机制
	return "admin" == user.LoginName
}

func (self defaultUserService) ToggleUserEnabled(user *entity.User) error {
	user.Enabled = !user.Enabled
	_, err := self.session.Id(user.Id).Cols("enabled").Update(user)
	return err
}

// TODO 缓存
func (self defaultUserService) GetUserLevel(user *entity.User) (level entity.UserLevel, ok bool) {
	var levels []entity.UserLevel

	self.session.Find(&levels)
	for _, level := range levels {
		log.Printf("%v", level)
		if matchLevel(user.Scores, level) {
			return level, true
		}
	}
	return level, false
}

func matchLevel(scores int, level entity.UserLevel) bool {
	if scores >= level.StartScores && scores <= level.EndScores {
		return true
	}
	if level.EndScores == 0 && scores >= level.StartScores {
		return true
	}
	return false
}

func (self defaultUserService) FindUserLevels() (levels []entity.UserLevel) {
	self.session.Find(&levels)
	return
}

func (self defaultUserService) FindAllUsersForPage(ps *PageSearcher) (page PageData) {
	ps.SearchKeyCall = func(session *xorm.Session) {
		session.Where("login_name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.User{})
	if err != nil {
		log.Println(err)
	}

	var users []entity.User
	err1 := ps.BuildQuerySession().Find(&users)
	if err1 != nil {
		log.Println(err1)
	}

	return NewPageData(total, users)
}

func (self defaultUserService) FindUserLoginLogs(userId int64) (llogs []entity.LoginLog) {
	_ = self.session.Where("user_id = ?", userId).Desc("id").Find(&llogs)
	return
}

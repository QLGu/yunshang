package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/deckarep/golang-set"
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
	Code      string
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
	return SessionUser{
		Id:        user.Id,
		Code:      user.Code,
		Email:     user.Email,
		LoginName: user.LoginName,
		RealName:  user.RealName,
		From:      from}
}

func NewUserService(db *xorm.Session) *UserService {
	return &UserService{db}
}

//////////////////////////////////////////////////////////////////////////////
// impls
type UserService struct {
	db *xorm.Session
}

func (self UserService) Total() int64 {
	ret, err := self.db.Count(UserTypeInstance)
	gotang.AssertNoError(err, "")

	return ret
}

func (self UserService) FindAllUsers() (users []entity.User) {
	self.db.Find(&users)
	return
}

func (self UserService) RegistUser(email, password string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = utils.Sha1(password)
	user.ActivationCode = utils.Uuid()
	user.Code = utils.Uuid()
	user.From = "Local"
	user.ActivationCodeCreatedAt = time.Now()

	_, err = self.db.Insert(&user)
	return
}

func (self UserService) ConnectUser(id string, providerName string, email string) (user entity.User, err error) {
	user.Email = email
	user.CryptedPassword = ""
	user.ActivationCode = ""
	user.LoginName = providerName + id
	user.From = providerName
	user.Code = utils.Uuid()
	user.Enabled = true

	_, err = self.db.Insert(&user)
	return
}

func (self UserService) ExistsUserByEmail(email string) bool {
	total, _ := self.db.Where("email=?", email).Count(UserTypeInstance)
	return total > 0
}

func (self UserService) Activate(email, code string) (user entity.User, err error) {
	var users []entity.User
	err = self.db.Where("email=? and activation_code=?", email, code).Find(&users)
	if err != nil {
		return
	}
	if len(users) > 0 {
		user = users[0]
		user.Enabled = true
		user.ActivationCode = ""
		self.db.Id(user.Id).Cols("enabled", "activation_code").Update(&user)
		return
	} else {
		err = errors.New("激活码不存在或已经失效！")
		return
	}
}

func (self UserService) CheckUser(login, password string) (user entity.User, ok bool) {
	ok, err := self.db.Where("(email=? or login_name=?) and crypted_password = ? and enabled=?", login, login, utils.Sha1(password), true).Get(&user)
	return user, ok && err == nil
}

func (self UserService) CheckUserByEmail(email string) (user entity.User, ok bool) {
	ok, err := self.db.Where("email=?", email).Get(&user)
	return user, ok && err == nil
}

func (self UserService) CheckUserByLoginName(loginName string) (user entity.User, ok bool) {
	ok, err := self.db.Where("login_name=?", loginName).Get(&user)
	return user, ok && err == nil
}

// 用户成功登录之后
func (self UserService) DoUserLogin(user *entity.User) error {
	// 更新最近一次登录时间
	user.LastSignAt = time.Now()
	_, err := self.db.Id(user.Id).Cols("last_sign_at").Update(user)

	date := user.LastSignAt.Format(gtime.ChinaDefaultDate)
	self.doUpdateLoginLogForEveryLogin(user.Id, date, user.LastSignAt)

	return err
}

// 更新登录日志： 策略： 每次登录都记录
func (self UserService) doUpdateLoginLogForEveryLogin(userId int64, date string, detailTime time.Time) (llog entity.LoginLog, new bool) {
	llog.Date = date
	llog.DetailTime = detailTime
	llog.UserId = userId
	self.db.Insert(&llog)
	new = true

	return
}

// 更新登录日志： 策略： 每天记录一条，并且更新到最近一次登录
func (self UserService) _doUpdateLoginLogForOneDay(userId int64, date string, detailTime time.Time) (llog entity.LoginLog, new bool) {
	exists, err := self.db.Where("date = ?", date).And("user_id = ?", userId).Get(&llog)
	gotang.AssertNoError(err, "")

	if exists {
		llog.DetailTime = detailTime
		self.db.Id(llog.Id).Cols("detail_time").Update(&llog)
		new = false
	} else {
		llog.Date = date
		llog.DetailTime = detailTime
		llog.UserId = userId
		self.db.Insert(&llog)
		new = true
	}
	return
}

func (self UserService) DoForgotPasswordApply(user *entity.User) error {
	user.PasswordResetCode = utils.Uuid()
	_, err := self.db.Id(user.Id).Update(user)
	return err
}

func (self UserService) ResetUserPassword(email, code string) (newPassword string, err error) {
	newPassword = utils.RandomString(6)
	var user entity.User
	ok, err := self.db.Where("email=? and password_reset_code = ?", email, code).Get(&user)
	if !ok {
		return "", errors.New("密码重置请求无效")
	}

	if err != nil {
		return "", err
	}

	user.CryptedPassword = utils.Sha1(newPassword)
	_, err = self.db.Id(user.Id).Cols("crypted_password").Update(&user)

	if err != nil {
		return "", err
	}
	return newPassword, err
}

func (self UserService) GetUserByLogin(login string) (user entity.User, ok bool) {
	ok, err := self.db.Where("email=? or login_name=?", login, login).Get(&user)
	return user, ok && err == nil
}

func (self UserService) DoChangePassword(user *entity.User, rawPassword string) error {
	user.CryptedPassword = utils.Sha1(rawPassword)
	_, err := self.db.Id(user.Id).Cols("crypted_password").Update(user)
	return err
}

func (self UserService) VerifyPassword(cryptedPassword, rawPassword string) bool {
	return cryptedPassword == utils.Sha1(rawPassword)
}

func (self UserService) GetUserById(id int64) (user entity.User, ok bool) {
	ok, err := self.db.Id(id).Get(&user)
	return user, ok && err == nil
}

func (self UserService) IsAdminUser(user *entity.User) bool {
	//TODO 改进判断机制
	return "admin" == user.LoginName
}

func (self UserService) ToggleUserEnabled(user *entity.User) error {
	user.Enabled = !user.Enabled
	_, err := self.db.Id(user.Id).Cols("enabled").Update(user)
	return err
}

func (self UserService) ToggleUserCertified(user *entity.User) error {
	user.Certified = !user.Certified
	_, err := self.db.Id(user.Id).Cols("certified").Update(user)
	return err
}

// TODO 缓存
func (self UserService) GetUserLevel(user *entity.User) (level entity.UserLevel, ok bool) {
	var levels []entity.UserLevel

	self.db.Find(&levels)
	for _, level := range levels {
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

func (self UserService) FindUserLevels() (levels []entity.UserLevel) {
	self.db.Find(&levels)
	return
}

func (self UserService) FindAllUsersForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("login_name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.User{})
	gotang.AssertNoError(err, "")

	var users []entity.User
	err1 := ps.BuildQuerySession().Find(&users)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, users, ps)
}

func (self UserService) CommentsForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.And("content like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Comment{})
	gotang.AssertNoError(err, "")

	var comments []entity.Comment
	err1 := ps.BuildQuerySession().Find(&comments)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, comments, ps)
}

// 列出用户登录的日志
func (self UserService) FindUserLoginLogs(userId int64) (llogs []entity.LoginLog) {
	_ = self.db.Where("user_id = ?", userId).Desc("id").Find(&llogs)
	return
}

// 批量计算用户的积分
// date： 要计算的日期
func (self UserService) ComputeUsersScores(date string) (err error) {
	const (
		INC_ONE       = 1
		INC_FOUR      = 4
		CONTINUE_DAYS = 7
	)

	dt, err := time.Parse(gtime.ChinaDefaultDate, date)
	gotang.AssertNoError(err, "")

	st := dt.AddDate(0, 0, -(CONTINUE_DAYS - 1))
	st_date := st.Format(gtime.ChinaDefaultDate)
	weekdates := []interface{}{date, st_date}
	for i := 1; i <= CONTINUE_DAYS-2; i++ {
		weekdates = append(weekdates, st.AddDate(0, 0, i).Format(gtime.ChinaDefaultDate))
	}

	//连续7天登录？
	isLoginWeek := func(userId int64) bool {
		var llogs []entity.LoginLog
		self.db.Where("user_id=?", userId).Cols("date").Find(&llogs)
		var reals []interface{}
		for _, v := range llogs {
			reals = append(reals, v.Date)
		}

		return len(mapset.NewSetFromSlice(weekdates).Difference(mapset.NewSetFromSlice(reals))) == 0
	}

	// 步骤1： 按登录计
	// 找出当天天有登录的用户
	// if 找出前6天都有登录的用户 + 4分
	// else + 1分
	err = self.db.Where("date = ?", date).Iterate(LoginLogTypeInstance, func(i int, bean interface{}) error {
		// 当天有登录记录
		llog := bean.(*entity.LoginLog)
		if isLoginWeek(llog.UserId) {
			return self.DoIncUserScores(llog.UserId, INC_FOUR)
		}
		return self.DoIncUserScores(llog.UserId, INC_ONE)
	})

	// TODO 步骤2：按评价计
	return
}

func (self UserService) DoIncUserScores(userId int64, inc int) error {
	user, ok := self.GetUserById(userId)
	if ok {
		user.Scores += inc
		_, err := self.db.Id(userId).Cols("scores").Update(&user)
		return err
	}
	return nil
}

func (self UserService) UpdateUserInfo(currUser *entity.User, user entity.User, userDetail entity.UserDetail) error {
	if len(currUser.LoginName) == 0 {
		currUser.LoginName = user.LoginName
	}
	if len(currUser.Email) == 0 {
		currUser.Email = user.Email
	}

	currUser.Gender = user.Gender
	currUser.MobilePhone = user.MobilePhone
	currUser.RealName = user.RealName

	if _, err := self.db.Id(currUser.Id).Cols("login_name", "email", "gender", "mobile_phone", "real_name").Update(currUser); err != nil {
		return err
	}

	var currUserDetails []entity.UserDetail
	err := self.db.Where("user_id = ?", currUser.Id).Find(&currUserDetails)
	if err != nil {
		return err
	}

	if len(currUserDetails) == 1 {
		userDetail.Id = currUserDetails[0].Id
		userDetail.UserId = currUser.Id
		if _, err := self.db.Id(userDetail.Id).Cols("work_kind", "id_number", "qq", "msn", "ali_wangwang", "birthday_year", "birthday_month",
			"birthday_day", "company_name", "company_type", "company_main_biz",
			"company_detail_biz", "company_website", "company_address", "company_zip_code", "company_fax",
			"company_phone", "company_province", "company_city", "company_area").Update(&userDetail); err != nil {
			return err
		}
	} else {
		userDetail.UserId = currUser.Id
		if _, err := self.db.Insert(&userDetail); err != nil {
			return err
		}
	}

	return nil
}

func (self UserService) GetUserDetailByUserId(userId int64) (userDetail entity.UserDetail, ok bool) {
	ok, _ = self.db.Where("user_id = ?", userId).Get(&userDetail)
	return
}

func (self UserService) FindUserDeliveryAddresses(userId int64) (das []entity.DeliveryAddress) {
	_ = self.db.Where("user_id=?", userId).Asc("id").Find(&das)
	return
}

func (self UserService) GetUserDeliveryAddress(userId int64, daId int64) (da entity.DeliveryAddress, ok bool) {
	ok, _ = self.db.Where("id=? and user_id=?", daId, userId).Get(&da)
	return
}

func (self UserService) SaveUserDeliveryAddress(da entity.DeliveryAddress) (id int64, err error) {
	if da.Id == 0 { //insert
		_, err = self.db.Insert(&da)
		id = da.Id
		return
	} else { // update
		currDa, ok := self.GetUserDeliveryAddress(da.UserId, da.Id)
		if ok {
			da.DataVersion = currDa.DataVersion
			_, err = self.db.Id(da.Id).Update(&da)

			if da.IsMain {
				das := self.FindUserDeliveryAddresses(da.UserId)
				for _, it := range das {
					it.IsMain = false
					self.db.Cols("is_main").Update(&it) // 批量更新以前的地址为no-main
				}
			}
			currDa, ok = self.GetUserDeliveryAddress(da.UserId, da.Id) // IsMain
			currDa.IsMain = da.IsMain
			_, err = self.db.Id(currDa.Id).Cols("is_main").Update(&currDa)

			return da.Id, err
		} else {
			return 0, fmt.Errorf("此收货地址不存在")
		}
	}
}

func (self UserService) GetUserDeliveryAddressTotal(userId int64) int64 {
	t, _ := self.db.Where("user_id=?", userId).Count(&entity.DeliveryAddress{})
	return t
}

func (self UserService) DeleteDeliveryAddress(userId, id int64) error {
	_, err := self.db.Where("id=? and user_id=?", id, userId).Delete(&entity.DeliveryAddress{})
	return err
}

func (self UserService) FindAllProductCollectsForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.And("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.ProductCollect{})
	gotang.AssertNoError(err, "")

	pcs := []entity.ProductCollect{}
	err1 := ps.BuildQuerySession().Find(&pcs)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, pcs, ps)
}

func (self UserService) CollectProduct(userId int64, productId int64) (err error) {
	count, _ := self.db.Where("product_id=? and user_id=?", productId, userId).Count(&entity.ProductCollect{})
	if count > 0 {
		return errors.New("您已经收藏了此产品！")
	}
	price := NewProductService(self.db).GetProductUnitPrice(productId)
	p := entity.ProductCollect{ProductId: productId, UserId: userId, Price: price}
	_, err = self.db.Insert(&p)
	return
}

func (self UserService) TotalUserCollects(userId int64) (count int64) {
	count, _ = self.db.Where("user_id=?", userId).Count(&entity.ProductCollect{})
	return
}

func (self UserService) FindUserInvoices(userId int64) (ps []entity.Invoice) {
	_ = self.db.Where("user_id=?", userId).Find(&ps)
	return
}

func (self UserService) GetUserInvoice(userId int64, id int64) (in entity.Invoice, ok bool) {
	ok, _ = self.db.Where("id=? and user_id=?", id, userId).Get(&in)
	return
}

func (self UserService) SaveUserInvoice(in entity.Invoice) (id int64, err error) {
	if in.Id == 0 { //insert
		_, err = self.db.Insert(&in)
		id = in.Id
		return
	} else { // update
		currIn, ok := self.GetUserInvoice(in.UserId, in.Id)
		if ok {
			_, err = self.db.Id(currIn.Id).Update(&in)
			return in.Id, err
		} else {
			return 0, fmt.Errorf("此发票信息不存在")
		}
	}
}

func (self UserService) GetUserInvoiceTotal(userId int64) int64 {
	t, _ := self.db.Where("user_id=?", userId).Count(&entity.Invoice{})
	return t
}

func (self UserService) DeleteInvoice(userId, id int64) error {
	_, err := self.db.Where("id=? and user_id=?", id, userId).Delete(&entity.Invoice{})
	return err
}

func (self UserService) GetUserDesc(userId int64) string {
	user, exists := self.GetUserById(userId)
	if !exists {
		return ""
	}
	userDetail, exists := self.GetUserDetailByUserId(user.Id)
	if !exists {
		return fmt.Sprintf("%s(%s)", user.LoginName, user.RealName)
	}
	return fmt.Sprintf("%s(%s), %s", user.LoginName, user.RealName, userDetail.CompanyName)
}

func (self UserService) CommentProducts(userId int64, ps []int64, scores int, content string) (err error) {
	user, exists := self.GetUserById(userId)
	gotang.Assert(exists, "用户不存在！")

	for _, p := range ps {
		product, _ := NewProductService(self.db).GetProductById(p)
		c := entity.Comment{
			UserId: userId, TargetId: p, TargetType: entity.CT_PRODUCT,
			TargetName: product.Name, Scores: scores, Content: content,
			UserName: user.DisplayName()}
		_, err = self.db.Insert(&c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self UserService) CommentNews(userId int64, newsId int64, name string, content string) (err error) {
	if userId != 0 {
		user, exists := self.GetUserById(userId)
		gotang.Assert(exists, "用户不存在！")
		name = user.DisplayName()
	}

	news, exists := NewNewsService(self.db).GetNewsById(newsId)
	gotang.Assert(exists, "新闻不存在！")

	c := entity.Comment{
		UserId: userId, TargetId: newsId, TargetType: entity.CT_ARTICLE,
		TargetName: news.Title, Content: content,
		UserName: name}
	_, err = self.db.Insert(&c)
	if err != nil {
		return err
	}

	return nil
}

func (self UserService) FindUserComments(userId int64) (ps []entity.Comment) {
	_ = self.db.Where("user_id=?", userId).Find(&ps)
	return
}

func (self UserService) FindUserCommentsForPage(userId int64, ps *PageSearcher) (page *PageData) {
	ps.FilterCall = func(session *xorm.Session) {
		session.And("user_id=?", userId)
	}
	total, err := ps.BuildCountSession().Count(&entity.Comment{})
	gotang.AssertNoError(err, "")

	var comments []entity.Comment
	err1 := ps.BuildQuerySession().Find(&comments)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, comments, ps)
}

func (self UserService) GetUserComment(userId, id int64) (c entity.Comment, exists bool) {
	exists, _ = self.db.Where("user_id=? and id=?", userId, id).Get(&c)
	return
}

func (self UserService) GetCommentById(id int64) (c entity.Comment, exists bool) {
	exists, _ = self.db.Where(" id=?", id).Get(&c)
	return
}

func (self UserService) DeleteComment(userId int64, id int64) (err error) {
	c, exists := self.GetUserComment(userId, id)
	if !exists {
		return errors.New("此评论不存在！")
	}
	_, err = self.db.Delete(c)
	return
}

func (self UserService) ToggleCommentEnabled(comment *entity.Comment) error {
	comment.Enabled = !comment.Enabled
	if comment.Enabled {
		comment.EnabledAt = time.Now()
	}
	_, err := self.db.Id(comment.Id).Cols("enabled", "enabled_at").Update(comment)
	if err != nil {
		return err
	}
	if comment.Enabled {
		FireEvent(EventObject{Name: EProductComment, UserId: comment.UserId, SourceId: comment.Id})
	}
	return err
}

func (self UserService) FindAllUserInquiries(userId int64) (ps []entity.Inquiry) {
	_ = self.db.Where("user_id=?", userId).Find(&ps)
	return
}

func (self UserService) DeleteInquiryByUser(userId int64, id int64) (err error) {
	var in entity.Inquiry
	exists, _ := self.db.Where("user_id=? and id=?", userId, id).Get(&in)
	if !exists {
		return errors.New("此询价不存在！")
	}
	_, err = self.db.Delete(&in)
	return
}

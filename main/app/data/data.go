package data

import (
	"strconv"
	"time"

	"github.com/itang/gotang"
	gtime "github.com/itang/gotang/time"
	. "github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/lunny/xorm"
)

func TryInitData(engine *xorm.Engine) {
	dataIniter := DataIniter{engine}
	if dataIniter.needInit() {
		dataIniter.initUsers()
		dataIniter.initUserLevels()

		//test
		dataIniter.initLoginLogsForTest()
	}
}

type DataIniter struct {
	engine *xorm.Engine
}

func (self DataIniter) needInit() bool {
	total, _ := self.engine.Where("login_name = ?", "admin").Count(&User{})
	return total == 0
}

func (self DataIniter) initUsers() {
	admin := User{
		Email: "livetang@qq.com", CryptedPassword: utils.Sha1("computer"),
		LoginName: "admin", Enabled: true, Scores: 70000, Gender: "male", RealName: "系统管理员", Code: utils.Uuid(),
		Certified: true,
	}
	test := User{
		Email: "test@test.com", CryptedPassword: utils.Sha1("computer"),
		LoginName: "test", Enabled: true, Gender: "female", RealName: "测试1", Code: utils.Uuid(),
	}

	users := []User{admin, test}
	for _, user := range users {
		_, err := self.engine.Insert(&user)
		gotang.AssertNoError(err, "")
	}
}

func (self DataIniter) initUserLevels() {
	levels := []UserLevel{
		UserLevel{Name: "童生", StartScores: 0, EndScores: 49},
		UserLevel{Name: "秀才", StartScores: 50, EndScores: 99},
		UserLevel{Name: "举人", StartScores: 100, EndScores: 299},
		UserLevel{Name: "进士", StartScores: 300, EndScores: 999},
		UserLevel{Name: "探花", StartScores: 1000, EndScores: 4999},
		UserLevel{Name: "榜眼", StartScores: 5000, EndScores: 14999},
		UserLevel{Name: "状元", StartScores: 15000, EndScores: 29999},
		UserLevel{Name: "大学士", StartScores: 30000, EndScores: 69999},
		UserLevel{Name: "翰林文圣", StartScores: 70000, EndScores: 0},
	}
	for index, level := range levels {
		level.Sort = index
		level.Code = strconv.Itoa(index)
		_, err := self.engine.Insert(&level)
		gotang.AssertNoError(err, "")
	}
}

func (self DataIniter) initLoginLogsForTest() {
	const NUM = 7
	dws := make([]time.Time, 0)
	for i := NUM; i > 0; i-- {
		dws = append(dws, time.Now().AddDate(0, 0, -i))
	}

	llogs := make([]LoginLog, 0)
	for _, dw := range dws {
		llogs = append(llogs, LoginLog{UserId: 1, Date: dw.Format(gtime.ChinaDefaultDate), DetailTime: dw})
	}

	llogs = append(llogs, LoginLog{UserId: 2, Date: dws[NUM-1].Format(gtime.ChinaDefaultDate), DetailTime: dws[NUM-1]})
	for _, llog := range llogs {
		_, err := self.engine.Insert(&llog)
		gotang.AssertNoError(err, "")
	}
}

package data

import (
	"log"
	"strconv"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/lunny/xorm"
)

func TryInitData(engine *xorm.Engine) {
	dataIniter := DataIniter{engine}
	if dataIniter.needInit() {
		log.Printf("init data")

		dataIniter.initUsers()
		dataIniter.initUserLevels()
	}
}

type DataIniter struct {
	engine *xorm.Engine
}

func (self DataIniter) needInit() bool {
	total, _ := self.engine.Where("login_name = ?", "admin").Count(&entity.User{})
	return total == 0
}

func (self DataIniter) initUsers() {
	admin := entity.User{Email: "livetang@qq.com", CryptedPassword: utils.Sha1("computer"), LoginName: "admin", Enabled: true, Scores: 70000}
	test := entity.User{Email: "test@test.com", CryptedPassword: utils.Sha1("computer"), LoginName: "test", Enabled: true}
	users := []entity.User{admin, test}
	for _, user := range users {
		_, err := self.engine.Insert(&user)
		gotang.AssertNoError(err)
	}
}

func (self DataIniter) initUserLevels() {
	t1 := entity.UserLevel{Name: "童生", StartScores: 0, EndScores: 49}
	t2 := entity.UserLevel{Name: "秀才", StartScores: 50, EndScores: 99}
	t3 := entity.UserLevel{Name: "举人", StartScores: 100, EndScores: 299}
	t4 := entity.UserLevel{Name: "进士", StartScores: 300, EndScores: 999}
	t5 := entity.UserLevel{Name: "探花", StartScores: 1000, EndScores: 4999}
	t6 := entity.UserLevel{Name: "榜眼", StartScores: 5000, EndScores: 14999}
	t7 := entity.UserLevel{Name: "状元", StartScores: 15000, EndScores: 29999}
	t8 := entity.UserLevel{Name: "大学士", StartScores: 30000, EndScores: 69999}
	t9 := entity.UserLevel{Name: "翰林文圣", StartScores: 70000, EndScores: 0}
	levels := []entity.UserLevel{t1, t2, t3, t4, t5, t6, t7, t8, t9}
	for index, level := range levels {
		level.Sort = index
		level.Code = strconv.Itoa(index)
		_, err := self.engine.Insert(&level)
		gotang.AssertNoError(err)
	}
}

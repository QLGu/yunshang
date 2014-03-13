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
		dataIniter.initProductCategories()
		dataIniter.initProviders()

		//test
		dataIniter.initLoginProductsForTest()
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
	_, err := self.engine.Insert(users)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initUserLevels() {
	levels := []UserLevel{
		{Name: "童生", StartScores: 0, EndScores: 49},
		{Name: "秀才", StartScores: 50, EndScores: 99},
		{Name: "举人", StartScores: 100, EndScores: 299},
		{Name: "进士", StartScores: 300, EndScores: 999},
		{Name: "探花", StartScores: 1000, EndScores: 4999},
		{Name: "榜眼", StartScores: 5000, EndScores: 14999},
		{Name: "状元", StartScores: 15000, EndScores: 29999},
		{Name: "大学士", StartScores: 30000, EndScores: 69999},
		{Name: "翰林文圣", StartScores: 70000, EndScores: 0},
	}
	for index, level := range levels {
		level.Sort = index
		level.Code = strconv.Itoa(index)
		_, err := self.engine.Insert(&level)
		gotang.AssertNoError(err, "")
	}
}

func (self DataIniter) initProductCategories() {
	cgs := []ProductCategory{
		{Name: "ITWChemtronics产品", Enabled: true},
		{Name: "防雷管", Enabled: true},
		{Name: "场效应管", Enabled: true},
		{Name: "AC-DC 开关电源控制芯片", Enabled: true},
		{Name: "自恢复保险丝", Enabled: true},
		{Name: "OKinternational产品", Enabled: true},
		{Name: "DYMAX戴马斯产品", Enabled: true},
	}

	_, err := self.engine.Insert(cgs)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initProviders() {
	ps := []Provider{
		{Name: "凯泰电子", ShortName: "凯泰", Enabled: true, Introduce: "一家专业的电子元器件配套供应商"},
		{Name: "东芝半导体股份有限公司", ShortName: "Toshiba", Enabled: true, Introduce: "东芝半导体股份有限公司"},
	}
	_, err := self.engine.Insert(ps)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initLoginProductsForTest() {
	ps := []Product{
		{Name: "松香型吸锡编带/吸锡线", ProviderId: 1, CategoryId: 1, Enabled: true, Code: 10001, Introduce: "松香型吸锡编带 松香型，可以最快，最安全的方式清除残留焊锡 1、无腐蚀、超纯的R型松香助焊剂 2、将PCB受到热损伤的危险降到最小 3、不会在PCB上留下离子污染"},
		{Name: "超级喷力全方位除尘剂", ProviderId: 1, CategoryId: 1, Enabled: true, Code: 10002, Introduce: "特大喷力；可以任何角度喷射；快速清洁任何物体可以任何角度喷射而不会有液体喷出，避免由此导致的敏感物体表面的冻坏或损坏"},
	}
	_, err := self.engine.Insert(ps)
	gotang.AssertNoError(err, "")
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
	_, err := self.engine.Insert(llogs)
	gotang.AssertNoError(err, "")
}

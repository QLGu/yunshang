package data

import (
	"strconv"
	"time"

	"github.com/itang/gotang"
	gtime "github.com/itang/gotang/time"
	"github.com/itang/yunshang/main/app/models"
	. "github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/db"
	"github.com/lunny/xorm"
)

func TryInitData(db *xorm.Engine) {
	dataIniter := DataIniter{db}
	if dataIniter.needInit() {
		dataIniter.initUsers()
		dataIniter.initUserLevels()
		dataIniter.initProviders()

		//test
		dataIniter.initLoginProductsForTest()
		dataIniter.initLoginLogsForTest()
	}
	dataIniter.initProductCategories()
	dataIniter.initApps()
	dataIniter.initPayments()
	dataIniter.initDefaultDas()
	dataIniter.initDefaultShippings()
	dataIniter.initNewsCategories()
}

type DataIniter struct {
	db *xorm.Engine
}

func (self DataIniter) needInit() bool {
	total, _ := self.db.Where("login_name = ?", "admin").Count(&User{})
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
	_, err := self.db.Insert(users)
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
		_, err := self.db.Insert(&level)
		gotang.AssertNoError(err, "")
	}
}

func (self DataIniter) initProductCategories() {
	count, _ := self.db.Count(&ProductCategory{})
	if count > 0 {
		return
	}

	cgs := []ProductCategory{
		{Name: "ITWChemtronics产品", Enabled: true},
		{Name: "防雷管", Enabled: true},
		{Name: "场效应管", Enabled: true},
		{Name: "AC-DC 开关电源控制芯片", Enabled: true},
		{Name: "自恢复保险丝", Enabled: true},
		{Name: "OKinternational产品", Enabled: true},
		{Name: "DYMAX戴马斯产品", Enabled: true},
		{Name: "吸锡编带,吸锡线", Enabled: true, ParentId: 1},
		{Name: "导电笔,清洁笔,润滑剂", Enabled: true, ParentId: 1},
		{Name: "三防漆,康富涂层", Enabled: true, ParentId: 1},
		{Name: "阻焊膜,阻焊蓝胶", Enabled: true, ParentId: 1},
		{Name: "防静电液", Enabled: true, ParentId: 1},
		{Name: "无铅产品", Enabled: true, ParentId: 1},
		{Name: "光纤清洁产品", Enabled: true, ParentId: 1},
		{Name: "除尘剂", Enabled: true, ParentId: 1},
		{Name: "冷冻液", Enabled: true, ParentId: 1},
		{Name: "擦拭布和湿擦拭棒", Enabled: true, ParentId: 1},
		{Name: "擦拭棒", Enabled: true, ParentId: 1},
		{Name: "清洁剂，除脂剂，润滑剂", Enabled: true, ParentId: 1},
		{Name: "电源防雷应用", Enabled: true, ParentId: 2},
		{Name: "安防防雷应用", Enabled: true, ParentId: 2},
		{Name: "LED照明防雷", Enabled: true, ParentId: 2},
		{Name: "排烟", Enabled: true, ParentId: 6},
		{Name: "手工焊接", Enabled: true, ParentId: 6},
		{Name: "对流返修", Enabled: true, ParentId: 6},
		{Name: "医疗级别UV胶", Enabled: true, ParentId: 7},
		{Name: "电子工业UV胶", Enabled: true, ParentId: 7},
		{Name: "一般工业UV胶", Enabled: true, ParentId: 7},
		{Name: "紫外光固化设备", Enabled: true, ParentId: 7},
	}

	db.Do(func(db *xorm.Session) error {
		s := models.NewProductService(db)
		for _, e := range cgs {
			_, err := s.SaveCategory(e)
			gotang.AssertNoError(err, "")
		}
		return nil
	})
}

func (self DataIniter) initProviders() {
	ps := []Provider{
		{Name: "凯泰电子", ShortName: "凯泰", Enabled: true, Introduce: "一家专业的电子元器件配套供应商"},
		{Name: "TEST-东芝半导体股份有限公司", ShortName: "Toshiba", Enabled: true, Introduce: "东芝半导体股份有限公司"},
	}
	_, err := self.db.Insert(ps)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initLoginProductsForTest() {
	ps := []Product{
		{Name: "TEST-松香型吸锡编带/吸锡线", ProviderId: 1, CategoryId: 1, Enabled: true, Code: 10001, EnabledAt: time.Now(), Introduce: "松香型吸锡编带 松香型，可以最快，最安全的方式清除残留焊锡 1、无腐蚀、超纯的R型松香助焊剂 2、将PCB受到热损伤的危险降到最小 3、不会在PCB上留下离子污染"},
		{Name: "TEST-超级喷力全方位除尘剂", ProviderId: 1, CategoryId: 1, Enabled: true, Code: 10002, EnabledAt: time.Now(), Introduce: "特大喷力；可以任何角度喷射；快速清洁任何物体可以任何角度喷射而不会有液体喷出，避免由此导致的敏感物体表面的冻坏或损坏"},
	}
	_, err := self.db.Insert(ps)
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
	_, err := self.db.Insert(llogs)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initApps() {
	count, _ := self.db.Where("type=?", ATSg).Count(&AppParams{})
	if count > 0 {
		return
	}

	_, err := self.db.Insert(&AppParams{Type: ATSg, Name: "标语", Value: "您好， 欢迎来到ICGOO，这里是国内领先的专业级电子元器件直购网站！"})
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initPayments() {
	count, _ := self.db.Count(&Payment{})
	if count > 0 {
		return
	}
	ps := []Payment{{Name: "支付宝", Description: "", Enabled: true},
		{Name: "网银在线", Description: "", Enabled: false},
		{Name: "银行转账", Description: "", Enabled: true},
	}
	_, err := self.db.Insert(ps)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initDefaultDas() {
	count, _ := self.db.Where("is_visit", true).Count(&DeliveryAddress{})
	if count > 0 {
		return
	}
	ps := []DeliveryAddress{{Name: "上门自提", IsVisit: true}}
	_, err := self.db.Insert(ps)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initDefaultShippings() {
	count, _ := self.db.Count(&Shipping{})
	if count > 0 {
		return
	}
	ps := []Shipping{
		{Name: "顺丰速运", Description: "", Enabled: true},
		{Name: "圆通速递", Description: "", Enabled: true},
		{Name: "申通快递", Description: "", Enabled: true},
		{Name: "满1000包邮", Description: "", Enabled: true},
		{Name: "上门自提", Description: "", Enabled: true},
	}
	_, err := self.db.Insert(ps)
	gotang.AssertNoError(err, "")
}

func (self DataIniter) initNewsCategories() {
	count, _ := self.db.Count(&NewsCategory{})
	if count > 0 {
		return
	}

	cgs := []NewsCategory{
		{Name: "公司动态", Enabled: true}, // 1
		{Name: "行业动态", Enabled: true},
		{Name: "技术方案", Enabled: true}, // 3
		{Name: "客服服务", Enabled: true},
		{Name: "购物指南", Enabled: true, ParentId: 4}, // 5
		{Name: "支付方式", Enabled: true, ParentId: 4},
		{Name: "配送服务", Enabled: true, ParentId: 4}, // 7
		{Name: "售后服务", Enabled: true, ParentId: 4},
		{Name: "帮助中心", Enabled: true, ParentId: 4}, // 9
		{Name: "关于我们", Enabled: true},              //10
	}

	articles := []News{
		{Title: "注册账号", CategoryId: 5, Enabled: true},
		{Title: "怎样询价", CategoryId: 5, Enabled: true},
		{Title: "怎样下订单", CategoryId: 5, Enabled: true},
		{Title: "常见问题", CategoryId: 5, Enabled: true},
		{Title: "支付方式", CategoryId: 6, Enabled: true},
		{Title: "发票制度", CategoryId: 6, Enabled: true},
		{Title: "退款制度", CategoryId: 6, Enabled: true},
		{Title: "运费收取标准", CategoryId: 7, Enabled: true},
		{Title: "配送时间及配送范围", CategoryId: 7, Enabled: true},
		{Title: "货物跟踪", CategoryId: 7, Enabled: true},
		{Title: "上门自提", CategoryId: 7, Enabled: true},
		{Title: "服务及质量保证承诺", CategoryId: 8, Enabled: true},
		{Title: "退换货政策", CategoryId: 8, Enabled: true},
		{Title: "售后常见问题解答", CategoryId: 8, Enabled: true},
		{Title: "退换货流程", CategoryId: 8, Enabled: true},
		{Title: "找回密码", CategoryId: 9, Enabled: true},
		{Title: "客户建议", CategoryId: 9, Enabled: true},
		{Title: "客服投诉", CategoryId: 9, Enabled: true},

		{Title: "关于凯特", CategoryId: 10, Enabled: true},
		{Title: "招贤纳士", CategoryId: 10, Enabled: true},
		{Title: "联系方式", CategoryId: 10, Enabled: true},
		{Title: "站点地图", CategoryId: 10, Enabled: true},
	}

	db.Do(func(db *xorm.Session) error {
		s := models.NewNewsService(db)
		for _, e := range cgs {
			_, err := s.SaveCategory(e)
			gotang.AssertNoError(err, "")
		}
		for _, e := range articles {
			_, err := s.SaveNews(e)
			gotang.AssertNoError(err, "")
		}
		return nil
	})
}

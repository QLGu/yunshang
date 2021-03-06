package entity

import (
	"strconv"
	"strings"
	"time"
)

// 用户等级
type UserLevel struct {
	Id          int64
	Sort        int
	Name        string
	Code        string `xorm:"unique"`
	StartScores int    `xorm:"int default 0"`
	EndScores   int    `xorm:"int default 0"`
	Description string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}

// 用户工作性质
type UserWorkKind struct {
	Id        int64
	Sort      int
	Name      string
	Code      string    `xorm:"unique"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// 用户
type User struct {
	Id              int64  `json:"id"`
	Code            string `xorm:"unique not null index" json:"code"`      // 内部编码
	LoginName       string `xorm:"unique index" json:"login_name"`         // 登录名
	CryptedPassword string `xorm:"varchar(64)" json:"-"`                   // 密码（加密）
	Email           string `xorm:"varchar(100) unique index" json:"email"` // 邮件账号
	RealName        string `xorm:"varchar(100)" json:"real_name"`          // 真实姓名
	Scores          int    `xorm:"int default 0" json:"scores"`            // 积分
	Level           string `xorm:"varchar(20)" json:"level"`               // 等级, 冗余字段
	From            string `json:"from"`                                   // 注册来源

	Gender      string    `xorm:"varchar(100)" json:"gender"`       // 性别, 取值 male|femal|“”
	MobilePhone string    `xorm:"varchar(100)" json:"mobile_phone"` // 手机号
	LastSignAt  time.Time `json:"last_sign_at"`                     // 最近一次登录时间
	Certified   bool      `json:"certified"`                        // 认证过的，有保证的

	CompanyId int `json:"company_id"` // 公司Id

	Enabled                 bool      `json:"enabled"` // 帐户有效
	AccountNonExpired       bool      // 帐户未过期
	CredentialsNonExpired   bool      // 凭据未过期
	AccountNonLocked        bool      // 帐户未锁定
	ActivationCode          string    // 激活码
	ActivationCodeCreatedAt time.Time // 激活码创建时间
	PasswordResetCode       string    //           ; 密码重置码

	InnerTags string `json:"inner_tags"` //内部标签

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`
}

func (e User) DisplayName() string {
	if len(e.RealName) > 0 {
		return e.RealName
	}
	return e.LoginName
}

func (e User) IsAdminUser() bool {
	return e.LoginName == "admin"
}

func (e User) HasRole(role string) bool {
	return strings.Contains(e.InnerTags, "#"+role)
}

// 登录日志
type LoginLog struct {
	Id         int64
	UserId     int64
	Date       string
	DetailTime time.Time
}

// job执行日志
type JobLog struct {
	Id   int64
	Name string
	Date string
}

// 用户详情
type UserDetail struct {
	Id       int64
	UserId   int64  // 关联用户
	WorkKind string // 工作性质

	IdNumber    string // 身份证号
	Qq          string // QQ号
	Msn         string // MSN号
	AliWangwang string // 阿里旺旺号

	BirthdayYear  string // 生日
	BirthdayMonth string
	BirthdayDay   string

	LocationId string // 位置Id

	CompanyName string
	CompanyType string // 公司类型， 取值 企业单位：1， 个体经营：2, 事业单位或社会团体：3

	CompanyMainBiz   string // 主要产品或服务
	CompanyDetailBiz string // 具体产品或服务
	CompanyAddress   string
	CompanyZipCode   string
	CompanyFax       string //传真
	CompanyPhone     string
	CompanyWebsite   string // 公司主页

	CompanyProvince string //省
	CompanyCity     string //城市
	CompanyArea     string // 地区
}

func (e UserDetail) CompanyFullAddress() string {
	var (
		pid = e.CompanyProvince
		cid = pid + e.CompanyCity
		did = cid + e.CompanyArea
	)
	return rd.GetById(pid) + rd.GetById(cid) + rd.GetById(did) + e.CompanyAddress
}

// 收货地址
type DeliveryAddress struct {
	Id      int64 `json:"id"`
	UserId  int64 `json:"user_id"`                       //关联用
	IsMain  bool  `xorm:"default false" json:"is_main"`  // 首要地址？
	IsVisit bool  `xorm:"default false" json:"is_visit"` //上门地址													   Is

	Name        string `json:"name"`         //地址命名
	Consignee   string `json:"consignee"`    // 收货人
	Province    string `json:"province"`     //省
	City        string `json:"city"`         //城市
	Area        string `json:"area"`         // 地区
	Street      string `json:"street"`       // 街道
	Address     string `json:"address"`      // 街道
	MobilePhone string `json:"mobile_phone"` //手机号码
	FixedPhone  string `json:"fixed_phone"`  // 固定号码
	Email       string `json:"email"`        // 邮箱

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`
}

func (e DeliveryAddress) FullDetailAddress() string {
	var (
		pid = e.Province
		cid = pid + e.City
		did = cid + e.Area

		pn = rd.GetById(pid)
		cn = rd.GetById(cid)
		dn = rd.GetById(did)
	)

	if cn == "县" {
		cn = ""
	}
	if cn == "市辖区" {
		cn = ""
	}
	if dn == "市辖区" {
		dn = ""
	}

	return pn + cn + dn + " " + e.Address
}

func (e DeliveryAddress) FullPhones() string {
	return e.MobilePhone + " " + e.FixedPhone
}

const (
	IN_COMMON = 1 //普通
	IN_SPEC   = 2 //专有
)

//发票信息
type Invoice struct {
	Id int64 `json:"id"`

	Type int `json:"type"`

	UserId int64 `json:"user_id"`

	CompanyName string `json:"company_name"` //公司名称

	CompanyAddress string `json:"company_address"` //公司地址
	CompanyPhone   string `json:"company_phone"`   //公司电话
	TaxRegNumber   string `json:"tax_reg_number"`  //税务登记号

	BankName    string `json:"bank_name"`    //开户银行
	BankAccount string `json:"bank_account"` //银行账号
	DaAddress   string `json:"da_address"`   //发票寄送地址
	DaZipCode   string `json:"da_zip_code"`  //发票寄送地址邮编

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (e Invoice) TypeDesc() string {
	switch e.Type {
	case IN_COMMON:
		return "增值税普通发票（不可抵扣）"
	case IN_SPEC:
		return "增值税专用发票（可抵扣）"
	default:
		return ""
	}
}

const (
	CT_PRODUCT = 1 //产品
	CT_ARTICLE = 2 //文章
)

var CTArray IntKVS = []IntKV{{CT_PRODUCT, "产品"}, {CT_ARTICLE, "文章"}}
var CTMap = CTArray.ToMap()

//评论
type Comment struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`

	Enabled    bool   `json:"enabled"`
	TargetId   int64  `json:"target_id"`
	TargetType int    `json:"target_type"`
	TargetName string `json:"target_name"` //评论对象，冗余
	Content    string `xorm:"varchar(100)" json:"content"`
	Scores     int    `json:"scores"` //打分

	UserId   int64  `json:"user_id"`   //评论人
	UserName string `json:"user_name"` //评论人名字

	EnabledAt time.Time `json:"enabled_at"`
	EnabledBy int64     `json:"enabled_by"`
}

func (e Comment) TypeDesc() string {
	v, ok := CTMap[strconv.Itoa(e.TargetType)]
	if !ok {
		return ""
	}
	return v
}

func (e Comment) CanDelete() bool {
	return !e.Enabled
}

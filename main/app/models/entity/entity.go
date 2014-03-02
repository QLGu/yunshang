package entity

import "time"

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
	Scores          int    `xorm:"int default 0" json:"-"`                 // 积分
	Level           string `xorm:"varchar(20)" json:"level"`               // 等级, 冗余字段
	From            string `json:"from"`                                   // 注册来源

	Gender      string    `xorm:"varchar(100)" json:"gender"`       // 性别
	MobilePhone string    `xorm:"varchar(100)" json:"mobile_phone"` // 手机号
	LastSignAt  time.Time `json:"last_sign_at"`                     // 最近一次登录时间

	CompanyId int `json:"company_id"` // 公司Id

	Enabled                 bool      `json:"enabled"` // 帐户有效
	AccountNonExpired       bool      // 帐户未过期
	CredentialsNonExpired   bool      // 凭据未过期
	AccountNonLocked        bool      // 帐户未锁定
	ActivationCode          string    // 激活码
	ActivationCodeCreatedAt time.Time // 激活码创建时间
	PasswordResetCode       string    //           ; 密码重置码

	CreatedAt   time.Time `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated" json:"updated_at"`
	DataVersion int       `xorm:"version '_version'"`
}

// 登录日志
type LoginLog struct {
	Id         int64
	UserId     int64
	Date       string
	DetailTime time.Time
}

// 位置
type Location struct {
	Id       int64
	Province string //省
	City     string //城市
	Area     string // 地区
	Street   string // 街道
}

// 用户详情
type UserDetail struct {
	Id       int64
	UserId   int    // 关联用户
	WorkKind string // 工作性质

	IdNumber    string // 身份证号
	ZipCode     string // 邮编
	fax         string //传真
	Qq          string // QQ号
	Msn         string // MSN号
	AliWangwang string // 阿里旺旺号

	Birthday struct { // 生日
		Year  string
		Month string
		Day   string
	} `xorm:"extends"`

	LocationId string // 位置Id
}

// 公司类型
type CompanyType struct {
	Id   int64
	Name string
	Code string `xorm:"unique"`
}

// 公司主要产品或服务
type CompanyMainBiz struct {
	Id   int64
	Name string
	Code string `xorm:"unique"`
}

// 公司具体产品或服务
type CompanyDetailBiz struct {
	Id   int64
	Name string
	Code string `xorm:"unique"`
}

// 公司
type Company struct {
	Id   int64
	Name string `xorm:"unique not null"`
	Type string

	MainBiz    string // 主要产品或服务
	DetailBiz  string // 具体产品或服务
	WebsiteUrl string // 公司主页
}

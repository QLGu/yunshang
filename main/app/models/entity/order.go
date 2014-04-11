package entity

import (
	"strconv"
	"time"
)

// 购物车
type Cart struct {
	Id int64 `json:"id"`

	UserId int64 `json:"user_id"`

	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

// 支付方式
type Payment struct {
	Id int64 `json:"id"`

	Name        string `json:"name"`
	Description string `xorm:"varchar(1000)" json:"description"`

	Enabled bool `json:"enabled"`
}

// 配送方式
type Shipping struct {
	Id int64 `json:"id"`

	Name        string `json:"name"`
	Description string `xorm:"varchar(1000)" json:"description"`

	Enabled bool `json:"enabled"`
}

const (
	OS_TEMP   = 1 //临时订单
	OS_SUBMIT = 2 //提交的订单
	OS_PAY    = 3 // 支付的订单
	OS_VERIFY = 4 // 审核过的
	OS_SHIP   = 5 //已发货
	OS_FINISH = 6 //完成的
	OS_CANEL  = 7 // 取消的订单
	OS_LOCK   = 8 // 锁定的订单
	OS_BACK   = 9 //退货的订单

	PM_ZF = 1 //支付宝
	PM_WY = 2 //网银
	PM_ZZ = 3 // 转账

	SP_SF = 1 //顺丰
	SP_YT = 2 //圆通
	SP_ST = 3 // 申通
	SP_ZT = 4 //自提
	SP_BY = 5 //包邮
)

var OSArray IntKVS = []IntKV{
	{OS_TEMP, "临时订单"},
	{OS_SUBMIT, "未支付"},
	{OS_PAY, "已支付"},
	{OS_VERIFY, "已确认/待发货"},
	{OS_SHIP, "已发货"},
	{OS_CANEL, "已取消"},
	{OS_LOCK, "已锁定"},
	{OS_FINISH, "已完成"},
	{OS_BACK, "已退货"}}
var OSMap = OSArray.ToMap()

var PMArray IntKVS = IntKVS{{PM_ZF, "支付宝"}, {PM_WY, "网上银行"}, {PM_ZZ, "银行转账"}}
var PMMap = PMArray.ToMap()

var SPArray IntKVS = IntKVS{{SP_SF, "顺丰速运"}, {SP_YT, "圆通快递"}, {SP_ST, "申通快递"}, {SP_ZT, "用户自提"}, {SP_BY, "包邮"}}
var SPMap = SPArray.ToMap()

// 订单
type Order struct {
	Id   int64 `json:"id"`
	Code int64 `json:"code"` //订单号

	UserId int64 `json:"user_id"` // 会员用户Id

	DaId       int64 `json:"da_id"`       // 收货地址
	PaymentId  int64 `json:"payment_id"`  //支付方式
	InvoiceId  int64 `json:"invoice_id"`  //发票信息
	ShippingId int64 `json:"shipping_id"` //配送方式

	Amount float64 `xorm:"Numeric" json:"amount"` //总计

	SubmitAt time.Time `json:"submit_at"` //提交时间
	PayAt    time.Time `json:"pay_at"`    // 付款时间
	CancelAt time.Time `json:"cancel_at"` // 取消时间
	VerifyAt time.Time `json:"verify_at"` //审核时间
	LockAt   time.Time `json:"lock_at"`   //锁定时间
	ShipAt   time.Time `json:"ship_at"`   //发货时间
	FinishAt time.Time `json:"finish_at"` //完成时间

	Status     int `json:"status"`      //状态
	PrevStatus int `json:"prev_status"` //锁定前的状态

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (e Order) IsSubmited() bool {
	return e.Status >= OS_SUBMIT
}

func (e Order) IsCancel() bool {
	return e.Status == OS_CANEL
}

func (e Order) IsLocked() bool {
	return e.Status == OS_LOCK
}

func (e Order) IsPayed() bool {
	return e.Status >= OS_PAY && !e.PayAt.IsZero()
}

func (e Order) CanCancel() bool {
	return e.Status == OS_SUBMIT
}

func (e Order) CanLock() bool {
	return e.Status != OS_LOCK && e.Status != OS_FINISH
}

func (e Order) NeedPay() bool {
	return e.Status == OS_SUBMIT
}

func (e Order) CanDelete() bool {
	return e.Status == OS_CANEL
}

func (e Order) CanReceipt() bool {
	return e.Status == OS_SHIP
}

func (e Order) CanComment() bool {
	return e.Status == OS_FINISH
}

func (e Order) StatusDesc() string {
	desc, exists := OSMap[strconv.Itoa(e.Status)]
	if !exists {
		return "未知"
	}
	return desc
}

func (e Order) PaymentDesc() string {
	desc, exists := PMMap[strconv.Itoa(int(e.PaymentId))]
	if !exists {
		return "未知"
	}
	return desc
}

func (e Order) IsZFPay() bool {
	return e.PaymentId == PM_ZF
}

func (e Order) IsWYPay() bool {
	return e.PaymentId == PM_WY
}
func (e Order) IsZZPay() bool {
	return e.PaymentId == PM_ZZ
}

//订单明细
type OrderDetail struct {
	Id        int64   `json:"id"`
	OrderId   int64   `json:"order_id"`
	ProductId int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `xorm:"Numeric" json:"price"` //单价
}

type OrderLog struct {
	Id      int64  `json:"id"`
	OrderId int64  `json:"order_id"`
	Message string `json:"message"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

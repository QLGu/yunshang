package models

import (
	"errors"
	"fmt"
	"time"

	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
)

func NewOrderService(session *xorm.Session) *OrderService {
	return &OrderService{session}
}

type OrderService struct {
	db *xorm.Session
}

func (self OrderService) UserCarts(userId int64) int64 {
	c, _ := self.db.Where("user_id=?", userId).Count(&entity.Cart{})
	return c
}

func (self OrderService) GetCartProduct(userId int64, productId int64) (cart entity.Cart, ok bool) {
	ok, _ = self.db.Where("user_id=? and product_id=?", userId, productId).Get(&cart)
	return
}

func (self OrderService) FindUserCarts(userId int64) (cs []entity.Cart) {
	_ = self.db.Where("user_id=?", userId).Asc("id").Find(&cs)
	return
}

func (self OrderService) FindUserCartProductPrices(userId int64) (ps []entity.Product) {
	carts := self.FindUserCarts(userId)
	ids, _ := From(carts).Select(func(t T) (T, error) { return t.(entity.Cart).ProductId, nil }).Results()
	ps = NewProductService(self.db).FindProductsByIds(asInt64Slice(ids))
	return
}

func (self OrderService) FindOrderProducts(userId int64, code int64) (ps []entity.Product) {
	items := self.GetOrderItems(userId, code)
	ids, _ := From(items).Select(func(t T) (T, error) { return t.(entity.OrderDetail).ProductId, nil }).Results()
	ps = NewProductService(self.db).FindProductsByIds(asInt64Slice(ids))
	return
}

func (self OrderService) AddProductToCart(userId int64, productId int64, quantity int) (err error) {
	c, exists := self.GetCartProduct(userId, productId)

	if exists { // update
		return errors.New("此产品已经存在于购物车！")
	}

	//new
	p, exists := NewProductService(self.db).GetProductById(productId)
	if !exists {
		return errors.New("此产品不存在！")
	}
	if !p.Enabled {
		return errors.New("此产品未上架！")
	}

	if quantity < p.MinNumberOfOrders {
		quantity = p.MinNumberOfOrders
	}
	c = entity.Cart{UserId: userId, ProductId: productId, Quantity: quantity}
	_, err = self.db.Insert(&c)
	return
}

func (self OrderService) DeleteCartProduct(id int64, userId int64) (err error) {
	var cart entity.Cart
	ok, _ := self.db.Where("id=? and user_id=?", id, userId).Get(&cart)
	if !ok {
		return errors.New("此购物车项不存在！")
	}

	_, err = self.db.Delete(cart)
	return
}

func (self OrderService) CleanCart(userId int64) (err error) {
	ps := self.FindUserCarts(userId)
	for _, p := range ps {
		_, err = self.db.Delete(&p)
		if err != nil {
			return err
		}
	}
	return
}

func (self OrderService) TotalUserCarts(userId int64) (count int64) {
	count, _ = self.db.Where("user_id=?", userId).Count(&entity.Cart{})
	return
}

func (self OrderService) MoveCartsToCollects(userId int64) (err error) {
	ps := self.FindUserCarts(userId)
	Users := NewUserService(self.db)
	for _, p := range ps {
		_ = Users.CollectProduct(userId, p.ProductId)
		_, err = self.db.Delete(&p)
		if err != nil {
			return err
		}
	}
	return
}

func (self OrderService) IncCartProductQuantity(id int64, userId int64, quantity int) (cart entity.Cart, err error) {
	ok, _ := self.db.Where("id=? and user_id=?", id, userId).Get(&cart)
	if !ok {
		err = errors.New("此购物车项不存在！")
		return
	}

	cart.Quantity = cart.Quantity + quantity
	p, exists := NewProductService(self.db).GetProductById(cart.ProductId)
	if !exists {
		err = errors.New("此购物车项对应产品不存在！")
		return
	}
	if cart.Quantity > p.StockNumber {
		err = errors.New("购买数量超过了库存！")
		return
	}
	if cart.Quantity < p.MinNumberOfOrders {
		err = errors.New("购买数量小于起订数量！")
		return
	}

	self.db.Id(cart.Id).Cols("quantity").Update(&cart)
	return
}

func (self OrderService) GetOrderById(id int64) (o entity.Order, exists bool) {
	exists, _ = self.db.Where("id=?", id).Get(&o)
	return
}

func (self OrderService) GetOrderByCode(code int64) (o entity.Order, exists bool) {
	exists, _ = self.db.Where("code=?", code).Get(&o)
	return
}

func (self OrderService) GetOrder(userId int64, code int64) (o entity.Order, exists bool) {
	exists, _ = self.db.Where("user_id=? and code=?", userId, code).Get(&o)
	return
}

func (self OrderService) UpdateOrderCode(o entity.Order) (entity.Order, error) {
	o.Code = o.Id + 10000
	_, err := self.db.Id(o.Id).Cols("code").Update(&o)
	return o, err
}

func (self OrderService) SaveTempOrder(userId int64, ps []entity.ParamsForNewOrder) (order entity.Order, err error) {
	var o entity.Order
	o.UserId = userId

	//计算总价
	for _, p := range ps {
		o.Amount += p.PrefPrice * float64(p.Quantity)
	}

	_, err = self.db.Insert(&o)
	o, exists := self.GetOrderById(o.Id)
	gotang.Assert(exists, "no exists")

	// 更新订单号
	order, err = self.UpdateOrderCode(o)

	//生成订单明细
	for _, p := range ps {
		od := entity.OrderDetail{}
		od.OrderId = order.Id
		od.ProductId = p.ProductId
		od.Price = p.PrefPrice
		od.Quantity = p.Quantity
		_, err := self.db.Insert(&od)
		gotang.AssertNoError(err, "")
	}
	return
}

func (self OrderService) GetOrderItems(userId int64, code int64) (ps []entity.OrderDetail) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		return
	}

	_ = self.db.Where("order_id=?", order.Id).Find(&ps)
	return
}

func (self OrderService) GetOrderItemsByAdmin(id int64) (ps []entity.OrderDetail) {
	_ = self.db.Where("order_id=?", id).Find(&ps)
	return
}

func (self OrderService) FindShippings(amount float64) (ps []entity.Shipping) {
	if amount >= 1000 {
		_ = self.db.Where("enabled=true and (name like '%包邮%' or name like '%自提%')").Find(&ps)
	} else {
		_ = self.db.Where("name not like '%包邮%' and enabled=true").Find(&ps)
	}
	return
}

func (self OrderService) FindAllShippings() (ps []entity.Shipping) {
	_ = self.db.Find(&ps)
	return
}

func (self OrderService) SaveShippings(ps []entity.Shipping) (err error) {
	c := 0
	for _, v := range ps {
		if v.Enabled {
			c += 1
		}
	}
	if c == 0 {
		return errors.New("至少选中一个配送方式")
	}

	for _, p := range ps {
		cp, exists := self.GetShippingById(p.Id)
		if !exists {
			continue
		}

		cp.Enabled = p.Enabled
		cp.Description = p.Description
		_, err := self.db.Id(cp.Id).Cols("description", "enabled").Update(&cp)
		if err != nil {
			return err
		}
	}
	return
}

func (self OrderService) GetShippingById(id int64) (s entity.Shipping, exists bool) {
	exists, err := self.db.Where("id=?", id).Get(&s)
	gotang.AssertNoError(err, "GetShippingById")

	return
}

func (self OrderService) FindAllPayments() (ps []entity.Payment) {
	_ = self.db.Find(&ps)
	return
}

func (self OrderService) FindAPayments() (ps []entity.Payment) {
	_ = self.db.Where("enabled=true").Find(&ps)
	return
}

func (self OrderService) SubmitOrder(userId int64, n entity.Order) (err error) {
	order, exists := self.GetOrder(userId, n.Code)
	if !exists {
		return errors.New("订单不存在!")
	}

	order.PaymentId = n.PaymentId
	order.InvoiceId = n.InvoiceId
	order.DaId = n.DaId
	order.Status = entity.OS_SUBMIT
	order.SubmitAt = time.Now()
	order.ShippingId = n.ShippingId
	_, err = self.db.Id(order.Id).Cols("payment_id", "da_id", "invoice_id", "status", "submit_at", "shipping_id").Update(&order)

	if err != nil {
		return
	}

	//清理相应购物车项
	ps := self.FindOrderProducts(userId, order.Code)
	for _, c := range ps {
		productId := c.Id
		var cart entity.Cart
		ok, _ := self.db.Where("product_id=?", productId).Get(&cart)
		if ok {
			self.db.Delete(&cart)
		}
	}

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单状态", Message: "用户提交订单"})

	return
}

func (self OrderService) FindSubmitedOrdersByUser(userId int64) (ps []entity.Order) {
	_ = self.db.Where("user_id=? and status >= ?", userId, entity.OS_SUBMIT).Desc("submit_at").Find(&ps)
	return
}

func (self OrderService) FindSubmitedOrders() (ps []entity.Order) {
	_ = self.db.Where(" status >= ?", entity.OS_SUBMIT).Desc("submit_at").Find(&ps)
	return
}

func (self OrderService) FindSubmitedOrdersForPage(ps *PageSearcher) (page *PageData) {
	ps.FilterCall = func(db *xorm.Session) {
		db.And(" status >= ?", entity.OS_SUBMIT)
	}
	ps.SearchKeyCall = func(db *xorm.Session) {
		//db.Where("login_name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Order{})
	gotang.AssertNoError(err, "")

	var orders []entity.Order
	err1 := ps.BuildQuerySession().Find(&orders)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, orders, ps)
}

func (self OrderService) FindAllInquiresForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("model like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Inquiry{})
	gotang.AssertNoError(err, "")

	var ins []entity.Inquiry
	err1 := ps.BuildQuerySession().Find(&ins)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, ins, ps)
}

func (self OrderService) FindNewOrders() (ps []entity.Order) {
	_ = self.db.Where(fmt.Sprintf("status in (%d, %d)", entity.OS_SUBMIT, entity.OS_PAY)).Desc("submit_at").Find(&ps)
	return
}

func (self OrderService) TotalNewOrders() (total int64) {
	total, err := self.db.Where(fmt.Sprintf("status in (%d, %d)", entity.OS_SUBMIT, entity.OS_PAY)).Count(&entity.Order{})
	gotang.AssertNoError(err, "")

	return
}

func (self OrderService) CancelOrderByUser(userId int64, code int64) (err error) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		return errors.New("订单不存在!")
	}

	if !order.CanCancel() {
		return errors.New("此订单不能取消， 有任何疑问请联系客服!")
	}

	order.Status = entity.OS_CANEL
	order.CancelAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "cancel_at").Update(&order)

	if err != nil {
		return
	}

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单状态", Message: "用户取消了此订单"})

	return
}

func (self OrderService) DeleteOrderByUser(userId int64, code int64) (err error) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		return errors.New("订单不存在!")
	}

	if !order.CanDelete() {
		return errors.New("此订单不能删除， 有任何疑问请联系客服!")
	}

	_, err = self.db.Delete(order)
	return
}

func (self OrderService) PayOrderByUserComment(userId int64, code int64, comment string) (err error) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		return errors.New("订单不存在!")
	}

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: comment})
	return
}

func (self OrderService) FindAllOrderLogsByUser(userId int64, code int64) (ps []entity.OrderLog, err error) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		err = errors.New("订单不存在!")
		return
	}
	_ = self.db.Where("order_id=?", order.Id).Find(&ps)
	return
}

func (self OrderService) TotalUserOrdersByStatus(userId int64, status int) (count int64) {
	count, _ = self.db.Where("user_id=? and status=?", userId, status).Count(&entity.Order{})
	return
}

func (self OrderService) GetDaForView(userId int64, id int64) (da entity.DeliveryAddress, ok bool) {
	ok, _ = self.db.Where("user_id=? and id=?", userId, id).Get(&da)
	return
}

func (self OrderService) GetInForView(userId int64, id int64) (in entity.Invoice, ok bool) {
	ok, _ = self.db.Where("user_id=? and id=?", userId, id).Get(&in)
	return
}

func (self OrderService) GetShippingForView(userId int64, id int64) (shipping entity.Shipping, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&shipping)
	return
}

func (self OrderService) GetPaymentForView(userId int64, id int64) (payment entity.Payment, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&payment)
	return
}

func (self OrderService) ToggleOrderLock(id int64) (err error) {
	order, exists := self.GetOrderById(id)
	if !exists {
		return errors.New("订单不存在!")
	}

	if order.IsLocked() {
		order.Status = order.PrevStatus
		order.LockAt = time.Time{} // nil time
		_, err = self.db.Id(order.Id).Cols("status", "lock_at").Update(&order)
		if err != nil {
			return
		}
		FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "商城已经解锁了此订单！"})
	} else if order.CanLock() {
		order.PrevStatus = order.Status
		order.Status = entity.OS_LOCK
		order.LockAt = time.Now()
		_, err = self.db.Id(order.Id).Cols("prev_status", "status", "lock_at").Update(&order)
		if err != nil {
			return
		}
		FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "商城已经锁定了此订单！"})
	} else {
		err = errors.New(("订单不能被锁定！"))
	}
	return
}

func (self OrderService) ChangeOrderPayed(id int64) (err error) {
	order, exists := self.GetOrderById(id)
	if !exists {
		return errors.New("订单不存在!")
	}

	order.Status = entity.OS_PAY
	order.PayAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "pay_at").Update(&order)
	if err != nil {
		return
	}

	//减库存
	err = NewProductService(self.db).SubProductStockNumbersByOrder(order)
	gotang.AssertNoError(err, "")

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "商城已确认收款"})
	FireEvent(EventObject{Name: EPay, UserId: order.UserId, Message: "商城已确认收款", Data: order.Amount})
	return
}

/*
func (self OrderService) PayOrderByAdminManual(code int64) (err error) {
	order, exists := self.GetOrderByCode(code)
	if !exists {
		return errors.New("订单不存在!")
	}

	if !order.NeedPay() {
		return errors.New("此订单不能支付， 有任何疑问请联系客服!")
	}

	order.Status = entity.OS_PAY
	order.PayAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "pay_at").Update(&order)
	if err != nil {
		return
	}

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Message: "商城已确认收款"})
	println("HERERERE!!!!")
	FireEvent(EventObject{Name: EPay, UserId:order.UserId, Message: "商城已确认收款", Data: order.Amount})

	return
}
*/

func (self OrderService) ChangeOrderVerify(id int64) (err error) {
	order, exists := self.GetOrderById(id)
	if !exists {
		return errors.New("订单不存在!")
	}

	order.Status = entity.OS_VERIFY
	order.VerifyAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "verify_at").Update(&order)
	if err != nil {
		return
	}
	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "商城已确认此订单待发货！"})
	return
}

func (self OrderService) ChangeOrderShiped(id int64) (err error) {
	order, exists := self.GetOrderById(id)
	if !exists {
		return errors.New("订单不存在!")
	}

	order.Status = entity.OS_SHIP
	order.ShipAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "ship_at").Update(&order)
	if err != nil {
		return
	}
	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "商城确认此订单已发货！"})
	return
}

func (self OrderService) ReceiptOrder(userId int64, code int64) (err error) {
	order, exists := self.GetOrder(userId, code)
	if !exists {
		return errors.New("订单不存在!")
	}

	order.Status = entity.OS_FINISH
	order.FinishAt = time.Now()
	_, err = self.db.Id(order.Id).Cols("status", "finish_at").Update(&order)
	if err != nil {
		return
	}

	FireEvent(EventObject{Name: EOrderLog, SourceId: order.Id, Title: "订单信息", Message: "用户确认此订单收货，订单完成！"})
	return
}

func (self OrderService) GetInquiryById(id int64) (in entity.Inquiry, exists bool) {
	exists, _ = self.db.Where("id=?", id).Get(&in)
	return
}

func (self OrderService) GetInquiryByUser(userId, id int64) (in entity.Inquiry, exists bool) {
	exists, _ = self.db.Where("user_id=? and id=?", userId, id).Get(&in)
	return
}

func (self OrderService) GetInquiryReplies(id int64) (ps []entity.InquiryReply) {
	_ = self.db.Where("inquiry_id=?", id).Find(&ps)
	return
}

func (self OrderService) SaveInquiryReply(reply entity.InquiryReply) (err error) {
	in, exists := self.GetInquiryById(reply.InquiryId)
	if !exists {
		return errors.New("询价不存在")
	}

	in.Replies += 1
	_, err = self.db.Id(in.Id).Cols("replies").Update(&in)

	if err != nil {
		return
	}

	_, err = self.db.Insert(reply)
	if err != nil {
		return
	}

	return
}

func (self OrderService) DeleteInquiryReply(id int64) (err error) {
	var reply entity.InquiryReply
	exists, _ := self.db.Where("id=?", id).Get(&reply)

	if !exists {
		return errors.New("回复不存在！")
	}

	in, exists := self.GetInquiryById(reply.InquiryId)
	if !exists {
		return errors.New("询价不存在")
	}

	in.Replies -= 1
	_, err = self.db.Id(in.Id).Cols("replies").Update(&in)

	_, err = self.db.Delete(&reply)

	return
}

func (self OrderService) TotalUserInquiries(userId int64) (count int64) {
	count, _ = self.db.Where("user_id=?", userId).Count(&entity.Inquiry{})
	return
}

func (self OrderService) TotalUserInquiryReplied(userId int64) (count int64) {
	count, _ = self.db.Where("user_id=? and replies > ? ", userId, 0).Count(&entity.Inquiry{})
	return
}

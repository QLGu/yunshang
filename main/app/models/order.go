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

	for _, p := range ps {
		o.Amount += p.PrefPrice * float64(p.Quantity)
	}

	_, err = self.db.Insert(&o)
	o, exists := self.GetOrderById(o.Id)
	gotang.Assert(exists, "no exists")

	order, err = self.UpdateOrderCode(o)

	//
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

func (self OrderService) FindShippings(amount float64) (ps []entity.Shipping) {
	if amount >= 1000 {
		_ = self.db.Where("name like '%包邮%' or name like '%自提%'").Find(&ps)
	} else {
		_ = self.db.Where("name not like '%包邮%'").Find(&ps)
	}
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
	_, err = self.db.Id(order.Id).Cols("payment_id", "da_id", "invoice_id", "status", "submit_at").Update(&order)

	//TODO
	//清理相应购物车项

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
	_, err = self.db.Id(order.Id).Cols("status").Update(&order)
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

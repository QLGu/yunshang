package models

import (
	"errors"
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

func (self OrderService) AddProductToCart(userId int64, productId int64, quantity int) (err error) {
	c, exists := self.GetCartProduct(userId, productId)

	if exists { // update
		return errors.New("此产品已经存在于购物车！")
		//c.Quantity = quantity
		//_, err = self.db.Id(c.Id).Cols("quantity").Update(&c)
		//return
	}

	//new
	p, exists := NewProductService(self.db).GetProductById(productId)
	if !exists {
		return errors.New("此产品不存在！")
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
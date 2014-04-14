package controllers

import (
	"bytes"

	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/itang/yunshang/modules/alipay"
	"github.com/kr/pretty"
	"github.com/revel/revel"
)

func (c User) Cart() revel.Result {
	carts := c.orderApi().FindUserCarts(c.forceSessionUserId())

	c.setChannel("index/cart")
	return c.Render(carts)
}

func (c User) CartData() revel.Result {
	carts := c.orderApi().FindUserCarts(c.forceSessionUserId())
	return c.RenderJson(Success("", carts))
}

func (c User) CartProductPrices() revel.Result {
	ps := c.orderApi().FindUserCartProductPrices(c.forceSessionUserId())
	return c.RenderJson(Success("", ps))
}

func (c User) OrderProducts(code string) revel.Result {
	ps := c.orderApi().FindOrderProducts(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) CartProductPrefPrice(productId int64, quantity int) revel.Result {
	ps := c.productApi().GetProductPrefPrice(productId, quantity)
	return c.RenderJson(Success("", ps))
}

func (c User) Prices() revel.Result {
	ins := c.userApi().FindAllUserInquiries(c.forceSessionUserId())

	c.setChannel("order/prices")
	return c.Render(ins)
}

//提交新订单（临时）
func (c User) DoNewOrder(ps []entity.ParamsForNewOrder) revel.Result {
	revel.INFO.Printf("%v", ps)
	filterPs := make([]entity.ParamsForNewOrder, 0)
	for _, p := range ps {
		if p.CartId != 0 {
			filterPs = append(filterPs, p)
		}
	}
	c.Validation.Required(len(filterPs) > 0).Message("请选择商品！")
	if ret := c.doValidate(User.Cart); ret != nil {
		return ret
	}

	order, err := c.orderApi().SaveTempOrder(c.forceSessionUserId(), filterPs)
	c.Validation.Required(err == nil).Message("保存订单出错， 请重试！")
	if ret := c.doValidate(User.Cart); ret != nil {
		c.setRollbackOnly()
		revel.ERROR.Printf("%v", err)
		return ret
	}

	return c.Redirect(routes.User.ConfirmOrder(order.Code))
}

//确认订单
func (c User) ConfirmOrder(code string) revel.Result {
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}
	if order.IsSubmited() {
		c.Flash.Error("订单不存在！")
		return c.Redirect(User.Cart)
	}

	shippings := c.orderApi().FindShippings(order.Amount)
	payments := c.orderApi().FindAPayments()

	c.setChannel("index/confirm_order")
	return c.Render(order, shippings, payments)
}

func (c User) OrderItems(code string) revel.Result {
	if isManager(c.Session) {
		ps := c.orderApi().GetOrderItemsByCode(code)
		return c.RenderJson(Success("", ps))
	}
	ps := c.orderApi().GetOrderItems(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) OrderData(code string) revel.Result {
	ps, _ := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) DoSubmitOrder(o entity.Order) revel.Result {
	revel.INFO.Printf("%# v", pretty.Formatter(o))
	c.Validation.Required(o.DaId != 0).Message("请选择收获地址")

	if ret := c.doValidate(routes.User.ConfirmOrder(o.Code)); ret != nil {
		return ret
	}

	err := c.orderApi().SubmitOrder(c.forceSessionUserId(), o)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(routes.User.ConfirmOrder(o.Code))
	}

	c.Flash.Success("提交订单成功!")
	return c.Redirect(routes.User.PayOrder(o.Code))
}

func (c User) PayOrder(code string) revel.Result {
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}

	c.setChannel("index/orders_pay")
	return c.Render(order)
}

func (c User) CancelOrder(code string) revel.Result {
	err := c.orderApi().CancelOrderByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) ViewOrder(code string) revel.Result {
	c.setChannel("index/orders_view")
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}

	orderBy := c.userApi().GetUserDesc(order.UserId)

	return c.Render(order, orderBy)
}

func (c User) DeleteOrder(code string) revel.Result {
	err := c.orderApi().DeleteOrderByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) PayOrderByUserComment(code string, comment string) revel.Result {
	err := c.orderApi().PayOrderByUserComment(c.forceSessionUserId(), code, comment)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) OrderLogsData(code string) revel.Result {
	ps, err := c.orderApi().FindAllOrderLogsByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ps))
}

//在线支付
func (c User) PayOnline(code string, bank string) revel.Result {
	order, err := c.orderApi().GetAndCheckOrderCanPay(code)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(routes.User.PayOrder(code))
	}

	config := models.GetAlipayConfig()
	var req alipay.Request
	if order.IsWYPay() {
		req = alipay.NewBankPayRequest()
		req.Defaultbank = bank
	} else if order.IsZFPay() {
		req = alipay.NewRequest()
	} else {
		c.Flash.Error("支付方式选择不支持！")
		return c.Redirect(routes.User.PayOrder(code))
	}

	req.OutTradeNo = order.Code
	req.TotalFee = order.Amount
	req.Subject = models.CacheSystem.GetConfig("site.basic.name") + "订单:" + req.OutTradeNo

	var buf bytes.Buffer
	alipay.NewPage(config, req, &buf)

	//return c.RenderText(buf.String())
	return c.RenderHtml(buf.String())
}

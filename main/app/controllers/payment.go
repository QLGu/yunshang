package controllers

import (
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/itang/yunshang/modules/alipay"
	"github.com/revel/revel"
)

type Payment struct {
	AppController
}

//异步回调
func (c Payment) AlipayNotify() revel.Result {
	resp, err := alipay.ParseResponse(models.GetAlipayConfig(), c.Controller.Request.Form)
	if err != nil {
		revel.WARN.Printf("%v", err)

		return c.RenderText(alipay.FailureFeedbackCode)
	}

	//已经支付成功？
	if c.orderApi().HasPayByUserFromAlipay(resp) {
		return c.RenderText(alipay.SuccessFeedbackCode)
	}

	// 写入支付结果
	if _, err := c.orderApi().DoPayByUserFromAlipay(resp); err != nil {
		//写入支付结果出错， 返回错误代码（让alipay重试）
		return c.RenderText(alipay.FailureFeedbackCode)
	}

	return c.RenderText(alipay.SuccessFeedbackCode)
}

//同步call
func (c Payment) AlipayReturn() revel.Result {
	resp, err := alipay.ParseResponse(models.GetAlipayConfig(), c.Params.Values)

	if err != nil {
		if err.Error() == "sign invalid" {
			return c.RenderText("请求不合法")
		}

		code := resp.OutTradeNo

		c.Flash.Error("支付未完成! 请重试。有任何疑问请联系客服!")
		return c.Redirect(routes.User.ViewOrder(code))
	}

	revel.INFO.Printf("%v: %v", "alipay支付完成", resp)

	if c.orderApi().HasPayByUserFromAlipay(resp) {
		revel.INFO.Printf("%v: %v", "此订单已经支付！!", resp)
		c.Flash.Error("此订单已经支付！!")
		code := resp.OutTradeNo
		return c.Redirect(routes.User.ViewOrder(code))
	}

	revel.INFO.Printf("%v", "写入支付结果...")
	if order, err := c.orderApi().DoPayByUserFromAlipay(resp); err != nil {
		c.Flash.Success("写入支付结果失败! 请重试。有任何疑问请联系客服!")
		return c.Redirect(routes.User.ViewOrder(order.Code))
	} else {
		revel.INFO.Printf("%v", "写入支付结果成功")

		c.Flash.Success("支付完成!")
		return c.Redirect(routes.User.ViewOrder(order.Code))
	}
}

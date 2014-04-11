package controllers

import (
	"fmt"

	"github.com/itang/yunshang/main/app/models"
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
		revel.INFO.Printf("%v", err)
		return c.RenderText("请求不合法")
	}

	revel.INFO.Printf("%v", resp)
	return c.RenderText(alipay.SuccessFeedbackCode)
}

//同步call
func (c Payment) AlipayReturn() revel.Result {
	resp, err := alipay.ParseResponse(models.GetAlipayConfig(), c.Controller.Request.Form)
	if err != nil {
		revel.INFO.Printf("%v", err)
		return c.RenderText("请求不合法")
	}

	if resp.TradeStatus == ("TRADE_FINISHED") || resp.TradeStatus == "TRADE_SUCCESS" {
		revel.INFO.Printf("%v", "支付完成")
	}

	revel.INFO.Printf("%v", resp)
	return c.RenderText(fmt.Sprintf("%v", reps))
}

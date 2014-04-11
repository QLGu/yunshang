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
		revel.WARN.Printf("%v", err)

		return c.RenderText(alipay.FailureFeedbackCode)
	}

	if err := c.orderApi().DoPayByUserFromAlipay(resp); err != nil {
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
		return c.RenderText("支付未完成！")
	}

	revel.INFO.Printf("%v: %v", "支付完成", resp)

	return c.RenderText(fmt.Sprintf("%s: %v", "支付完成", resp))
}

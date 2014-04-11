package models

import (
	"log"
	"net/url"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/modules/alipay"
)

var alipayConfig alipay.Config

func InitAlipayConfig() {
	log.Println("InitAlipayConfig")

	alipayConfig.Partner = CacheSystem.GetConfig("site.alipay.partner")
	alipayConfig.Key = CacheSystem.GetConfig("site.alipay.key")

	alipayConfig.ReturnUrl = checkUrlWitchPanic(CacheSystem.UrlWithHost(CacheSystem.GetConfig("site.alipay.return_url")))
	alipayConfig.NotifyUrl = checkUrlWitchPanic(CacheSystem.UrlWithHost(CacheSystem.GetConfig("site.alipay.notify_url")))

	alipayConfig.PaymentType = CacheSystem.GetConfig("site.alipay.payment_type")
	alipayConfig.SellerEmail = CacheSystem.GetConfig("site.alipay.seller_email")
	alipayConfig.Service = CacheSystem.GetConfig("site.alipay.service")
}

func GetAlipayConfig() alipay.Config {
	return alipayConfig
}

func checkUrlWitchPanic(rawurl string) string {
	_, err := url.Parse(rawurl)
	gotang.AssertNoError(err, "URL:"+rawurl+"不合法")
	return rawurl
}

package controllers

import (
	"strings"

	"github.com/itang/iptaobao"
	"github.com/itang/yunshang/main/app"
	"github.com/revel/revel"
)

// 应用主控制器
type App struct {
	AppController
}

// 应用主页
func (c App) Index() revel.Result {
	version := app.Version
	ip := c.getRemoteIp()

	from := ""
	if !strings.HasPrefix(ip, "127") {
		r, err := iptaobao.GetIpInfo(ip)
		if err != nil {
			revel.INFO.Printf("%v, %v", ip, err)
		} else {
			from = r.Region + " " + r.City
		}
	}

	products := c.productApi().FindAllAvailableProducts()

	categories := c.productApi().FindAvailableTopCategories()

	return c.Render(version, ip, from, products, categories)
}

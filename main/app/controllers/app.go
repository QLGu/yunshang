package controllers

import (
	"strings"

	"github.com/itang/iptaobao"
	"github.com/itang/yunshang/main/app"
	"github.com/robfig/revel"
)

type App struct {
	AppController
}

func (c App) Index() revel.Result {
	version := app.Version
	ip := c.getRemoteIp()
	//ip = "27.46.125.49"
	revel.INFO.Printf("remoteAddr: %v", ip)

	from := ""
	if !strings.HasPrefix(ip, "127") {
		r, err := iptaobao.GetIpInfo(ip)
		if err != nil {
			revel.INFO.Printf("%v, %v", ip, err)
		} else {
			from = r.Region + " " + r.City
		}
	}

	return c.Render(version, ip, from)
}

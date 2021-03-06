package controllers

import (
	"fmt"
	"os"

	"github.com/itang/yunshang/main/app"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/revel/revel"
)

// 应用主控制器
type App struct {
	AppController
}

// 应用主页
func (c App) Index() revel.Result {
	c.setChannel("index/")
	return c.Render()
}

func (c App) AdImagesData() revel.Result {
	images := c.appApi().FindAdImages()
	return c.RenderJson(Success("", images))

}

func (c App) AdImage(file string) revel.Result {
	targetFile, err := c.appApi().GetAdImageFile(file)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

func (c App) HotKeywordsData() revel.Result {
	keywords := c.appApi().FindHotKeywords()
	return c.RenderJson(Success("", keywords))
}

func (c App) NewInquiry(q string) revel.Result {
	c.setChannel("index/inquiry")

	user, _ := c.currUser()
	userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)

	return c.Render(q, user, userDetail)
}

func (c App) DoNewInquiry(i entity.Inquiry) revel.Result {
	c.Validation.Required(i.Model).Message("请填写询价型号")
	c.Validation.Required(i.Quantity).Message("请填写询价数量")
	c.Validation.Required(i.Contact).Message("请填写联系人")
	c.Validation.Required(i.Phone).Message("请填写联系电话")

	if ret := c.doValidate(App.NewInquiry); ret != nil {
		return ret
	}

	if c.isLogined() {
		i.UserId = c.forceSessionUserId()
	}

	err := c.appApi().SaveInquiry(i)
	if err != nil {
		c.FlashParams()
		c.Flash.Error("询价出错， 请重试!")
		return c.Redirect(App.NewInquiry)
	}

	c.setChannel("index/inquiry")
	return c.Render()
}

func (c App) Version() revel.Result {
	return c.RenderText(app.Version)
}

func (c App) ProcessInfo() revel.Result {
	d := struct {
		Pid  int
		Ppid int
	}{os.Getpid(), os.Getppid()}
	return c.RenderJson(Success("进程信息", d))
}

func (c App) ProcessInfoLine() revel.Result {
	t := fmt.Sprintf("%d %d", os.Getpid(), os.Getppid())
	return c.RenderText(t)
}

func (c App) Weixin() revel.Result {
	return c.Render()
}

func (c App) NewFeedback(subject string) revel.Result {
	c.setChannel("index/feeback")

	user, _ := c.currUser()
	userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)

	return c.Render(subject, user, userDetail)
}

func (c App) DoNewFeedback(i entity.Feedback) revel.Result {
	c.Validation.Required(i.Content).Message("请填写内容")
	c.Validation.Required(i.Phone).Message("请填写联系电话")

	if ret := c.doValidate(App.NewFeedback); ret != nil {
		return ret
	}

	err := c.appApi().SaveFeedback(i)
	if err != nil {
		c.setRollbackOnly()
		c.FlashParams()
		c.Flash.Error("填写信息反馈出错， 请重试!")
		return c.Redirect(App.NewFeedback)
	}

	c.Flash.Success("提交信息反馈完成!")
	return c.Redirect(App.NewFeedback)
}

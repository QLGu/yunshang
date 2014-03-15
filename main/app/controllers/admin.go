package controllers

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/itang/yunshang/modules/mail"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 管理端相关Actions
type Admin struct {
	AdminController
}

// 管理端主页
func (c Admin) Index() revel.Result {
	_, ok := c.Session["locked"]
	if ok {
		c.Redirect(Admin.Lock)
	}

	userTotal := c.userApi().Total()
	orderTotal := 0 // TODO order total

	c.setChannel("/")
	return c.Render(userTotal, orderTotal)
}

// 锁屏
func (c Admin) Lock() revel.Result {
	c.Session["locked"] = "true"
	return c.Render()
}

// 解锁屏
func (c Admin) UnLock(password string) revel.Result {
	delete(c.Session, "locked")
	return c.Redirect(Admin.Index)
}

// 用户列表
func (c Admin) Users() revel.Result {
	c.setChannel("users/users")
	return c.Render()
}

// 用户列表数据
func (c Admin) UsersData(filter_status string, filter_certified string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}

		switch filter_certified {
		case "true":
			session.And("certified=?", true)
		case "false":
			session.And("certified=?", false)
		}
	})
	page := c.userApi().FindAllUsersForPage(ps)
	return c.renderDataTableJson(page)
}

// 重置用户密码
func (c Admin) ResetUserPassword(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(c.errorResposne("admin用户的状态不能通过此入口修改", nil))
	}

	newPassword := utils.RandomString(6)
	err := c.userApi().DoChangePassword(&user, newPassword)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	}

	err = gotang.DoIOWithTimeout(func() error {
		return mail.SendHtml("重置密码邮件",
			utils.RenderTemplate("Passport/ResetPasswordResultTemplate.html",
				struct {
					NewPassword string
				}{newPassword}),
			user.Email)
	}, time.Second*30)
	if err != nil {
		panic(err)
	}

	return c.RenderJson(c.successResposne("重置用户密码成功并新密码已经通过告知邮件用户", newPassword))
}

// 激活/禁用用户
func (c Admin) ToggleUserEnabled(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(c.errorResposne("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userApi().ToggleUserEnabled(&user)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		if user.Enabled {
			return c.RenderJson(c.successResposne("激活用户成功！", nil))
		}
		return c.RenderJson(c.successResposne("禁用用户成功！", nil))
	}
}

// 认证用户
func (c Admin) ToggleUserCertified(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(c.errorResposne("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userApi().ToggleUserCertified(&user)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		if user.Certified {
			return c.RenderJson(c.successResposne("设置用户认证成功！", nil))
		}
		return c.RenderJson(c.successResposne("取消用户认证成功！", nil))
	}
}

// 显示用户登录日志
func (c Admin) ShowUserLoginLogs(id int64) revel.Result {
	loginLogs := c.userApi().FindUserLoginLogs(id)
	return c.Render(loginLogs)
}

// 显示用户信息
func (c Admin) ShowUserInfos(id int64) revel.Result {
	user, _ := c.userApi().GetUserById(id)
	userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)
	userDas := c.userApi().FindUserDeliveryAddresses(user.Id)

	return c.Render(user, userDetail, userDas)
}

///////////////////////////////////////////////////////////////
// Products

// 产品列表
func (c Admin) Products() revel.Result {
	c.setChannel("products/products")
	return c.Render()
}

// 产品列表数据
func (c Admin) ProductsData(filter_status string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
	})
	page := c.productApi().FindAllProductsForPage(ps)
	return c.renderDataTableJson(page)
}

func (c Admin) NewProduct(id int64) revel.Result {
	var p entity.Product
	if id == 0 { // new
		//p = entity.Product{}
	} else { //edit
		p, _ = c.productApi().GetProductById(id)
	}
	return c.Render(p)
}

func (c Admin) DoNewProduct(p entity.Product) revel.Result {
	c.Validation.Required(p.Name).Message("请填写名称")

	if ret := c.doValidate(fmt.Sprintf("/admin/products/new?id=%d", p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveProduct(p)
	if err != nil {
		c.Flash.Error("保存产品失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存产品成功！")
	}

	return c.Redirect(fmt.Sprintf("/admin/products/new?id=%d", id))
}

func (c Admin) ToggleProductEnabled(id int64) revel.Result {
	api := c.productApi()
	p, ok := api.GetProductById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("产品不存在", nil))
	}

	err := api.ToggleProductEnabled(&p)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(c.successResposne("上架成功！", nil))
		}
		return c.RenderJson(c.successResposne("下架成功！", nil))
	}
}

func (c Admin) UploadProductImage(id int64, t int) revel.Result {
	var dir string
	var ct string
	if t == entity.PTScheDiag {
		dir = "data/products/sd/"
		ct = "thumbnail"
	} else if t == entity.PTPics {
		dir = "data/products/pics/"
		ct = "fit"
	} else {
		return c.RenderJson(c.errorResposne("上传失败！ 类型不对", nil))
	}

	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			p := entity.ProductParams{Type: t, Name: fileHeader.Filename, ProductId: id}
			e, err := c.XOrmSession.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to := utils.Uuid() + ".jpg"
			p.Value = to
			c.XOrmSession.Id(p.Id).Cols("value").Update(&p)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, ct, 200, 200)
			gotang.AssertNoError(err, "生成图片出错！")
		}
	}

	return c.RenderJson(c.successResposne("上传成功！", nil))
}

func (c Admin) UploadProductMaterial(id int64) revel.Result {
	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			to := ""
			p := entity.ProductParams{Type: entity.PTMaterial, Name: fileHeader.Filename, Value: to, ProductId: id}
			e, err := c.XOrmSession.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to = fmt.Sprintf("%d-%s", p.Id, fileHeader.Filename)
			p.Value = to
			c.XOrmSession.Id(p.Id).Cols("value").Update(&p)

			out, err := os.Create("data/products/m/" + to)
			gotang.AssertNoError(err, `os.Create`)

			in, err := fileHeader.Open()
			gotang.AssertNoError(err, `fileHeader.Open()`)

			io.Copy(out, in)

			out.Close()
			in.Close()
		}
	}

	return c.RenderJson(c.successResposne("上传成功！", nil))
}

func (c Admin) DeleteSdImage(id int64) revel.Result {
	var it entity.ProductParams
	_, _ = c.XOrmSession.Where("id=?", id).Get(&it)
	c.XOrmSession.Delete(&it)
	return c.RenderJson(c.successResposne("删除成功！", ""))
}

func (c Admin) DeleteImagePic(id int64) revel.Result {
	var it entity.ProductParams
	_, _ = c.XOrmSession.Where("id=?", id).Get(&it)
	c.XOrmSession.Delete(&it)
	return c.RenderJson(c.successResposne("删除成功！", ""))
}

func (c Admin) DeleteMFile(id int64) revel.Result {
	var it entity.ProductParams
	_, _ = c.XOrmSession.Where("id=?", id).Get(&it)
	c.XOrmSession.Delete(&it)
	return c.RenderJson(c.successResposne("删除成功！", ""))
}

func (c Admin) DoSaveProductSpec(productId int64, id int64, name string, value string) revel.Result {
	pp := entity.ProductParams{ProductId: productId, Id: id, Name: name, Value: value, Type: entity.PTSpec}
	if id == 0 { //new
		c.XOrmSession.Insert(&pp)
	} else { //update
		c.XOrmSession.Id(id).Update(&pp)
	}
	return c.RenderJson(c.successResposne("操作完成！", ""))
}

func (c Admin) DeleteSpec(id int64) revel.Result {
	var it entity.ProductParams
	_, _ = c.XOrmSession.Where("id=?", id).Get(&it)
	c.XOrmSession.Delete(&it)
	return c.RenderJson(c.successResposne("删除成功！", ""))
}

/////////////////////////////////////////////////////

func (c Admin) Categories() revel.Result {
	c.setChannel("products/categories")
	return c.Render()
}

func (c Admin) CategoriesData() revel.Result {
	page := c.productApi().FindAllCategoriesForPage(c.pageSearcher())
	return c.renderDataTableJson(page)
}

func (c Admin) NewCategory(id int64) revel.Result {
	var p entity.ProductCategory
	if id == 0 { // new
		//p = entity.Provider{}
	} else { //edit
		p, _ = c.productApi().GetCategoryById(id)
	}
	return c.Render(p)
}

func (c Admin) DoNewCategory(p entity.ProductCategory) revel.Result {
	c.Validation.Required(p.Name).Message("请填写名称")

	if ret := c.doValidate(fmt.Sprintf("/admin/categories/new?id=%d", p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveCategory(p)
	if err != nil {
		c.Flash.Error("保存分类失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存分类成功！")
	}

	return c.Redirect(fmt.Sprintf("/admin/categories/new?id=%d", id))
}

func (c Admin) ToggleCategoryEnabled(id int64) revel.Result {
	api := c.productApi()
	p, ok := api.GetCategoryById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("分类不存在", nil))
	}

	err := api.ToggleCategoryEnabled(&p)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(c.successResposne("激活成功！", nil))
		}
		return c.RenderJson(c.successResposne("禁用成功！", nil))
	}
}

//////////////////////////////////////////////////////////////////
// Providers

func (c Admin) Providers() revel.Result {
	c.setChannel("providers/providers")
	return c.Render()
}

func (c Admin) ProvidersData() revel.Result {
	page := c.productApi().FindAllProvidersForPage(c.pageSearcher())
	return c.renderDataTableJson(page)
}

func (c Admin) NewProvider(id int64) revel.Result {
	var p entity.Provider
	if id == 0 { // new
		//p = entity.Provider{}
	} else { //edit
		p, _ = c.productApi().GetProviderById(id)
	}
	return c.Render(p)
}

func (c Admin) DoNewProvider(p entity.Provider) revel.Result {
	c.Validation.Required(p.Name).Message("请填写名称")

	if ret := c.doValidate(fmt.Sprintf("/admin/providers/new?id=%d", p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveProvider(p)
	if err != nil {
		c.Flash.Error("保存制造商失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存制造商成功！")
	}

	return c.Redirect(fmt.Sprintf("/admin/providers/new?id=%d", id))
}

func (c Admin) ToggleProviderEnabled(id int64) revel.Result {
	api := c.productApi()
	p, ok := api.GetProviderById(id)
	if !ok {
		return c.RenderJson(c.errorResposne("制造商不存在", nil))
	}

	err := api.ToggleProviderEnabled(&p)
	if err != nil {
		return c.RenderJson(c.errorResposne(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(c.successResposne("激活成功！", nil))
		}
		return c.RenderJson(c.successResposne("禁用成功！", nil))
	}
}

func (c Admin) DeleteProvider(id int64) revel.Result {
	_ = c.productApi().DeleteProvider(id)
	return c.RenderJson(c.successResposne("删除成功！", nil))
}

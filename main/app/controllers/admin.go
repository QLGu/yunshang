package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	return c.renderDTJson(page)
}

// 重置用户密码
func (c Admin) ResetUserPassword(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(Error("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(Error("admin用户的状态不能通过此入口修改", nil))
	}

	newPassword := utils.RandomString(6)
	err := c.userApi().DoChangePassword(&user, newPassword)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
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

	return c.RenderJson(Success("重置用户密码成功并新密码已经通过告知邮件用户", newPassword))
}

// 激活/禁用用户
func (c Admin) ToggleUserEnabled(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(Error("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(Error("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userApi().ToggleUserEnabled(&user)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if user.Enabled {
			return c.RenderJson(Success("激活用户成功！", nil))
		}
		return c.RenderJson(Success("禁用用户成功！", nil))
	}
}

// 认证用户
func (c Admin) ToggleUserCertified(id int64) revel.Result {
	user, ok := c.userApi().GetUserById(id)
	if !ok {
		return c.RenderJson(Error("用户不存在", nil))
	}

	if c.userApi().IsAdminUser(&user) {
		return c.RenderJson(Error("admin用户的状态不能通过此入口修改", nil))
	}

	err := c.userApi().ToggleUserCertified(&user)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if user.Certified {
			return c.RenderJson(Success("设置用户认证成功！", nil))
		}
		return c.RenderJson(Success("取消用户认证成功！", nil))
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
	return c.renderDTJson(page)
}

func (c Admin) NewProduct(id int64) revel.Result {
	var (
		p         entity.Product
		detail    = ""
		stockLogs []entity.ProductStockLog
	)

	if id == 0 { // new

	} else { //edit
		p, _ = c.productApi().GetProductById(id)
		detail, _ = c.productApi().GetProductDetail(p.Id)
		stockLogs = c.productApi().FindAllProductStockLogs(p.Id)
	}
	return c.Render(p, detail, stockLogs)
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
		return c.RenderJson(Error("产品不存在", nil))
	}

	err := api.ToggleProductEnabled(&p)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(Success("上架成功！", nil))
		}
		return c.RenderJson(Success("下架成功！", nil))
	}
}

func (c Admin) UploadProductImage(id int64, t int) revel.Result {
	var (
		dir, ct string
		count   int
	)

	if t == entity.PTScheDiag {
		dir = "data/products/sd/"
		ct = "thumbnail"
	} else if t == entity.PTPics {
		dir = "data/products/pics/"
		ct = "fit"
	} else {
		return c.RenderJson(Error("上传失败！ 类型不对", nil))
	}

	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			p := entity.ProductParams{Type: t, Name: fileHeader.Filename, ProductId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to := utils.Uuid() + ".jpg"
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, ct, 200, 200)
			gotang.AssertNoError(err, "生成图片出错！")

			count += 1
		}
	}

	if count == 0 {
		return c.RenderJson(Error("请选择要上传的图片", nil))
	}

	return c.RenderJson(Success("上传成功！", nil))
}

func (c Admin) UploadProductImageForUEditor(id int64) revel.Result {
	dir := "data/products/pics/"
	ct := "fit"
	t := entity.PTPics

	var Original = ""
	var Url = ""
	var Title = ""
	var State = ""
	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			p := entity.ProductParams{Type: t, Name: fileHeader.Filename, ProductId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to := utils.Uuid() + ".jpg"
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, ct, 200, 200)
			gotang.AssertNoError(err, "生成图片出错！")

			Original = fileHeader.Filename
			Title = Original
			State = "SUCCESS"
			Url = "?file=" + to
		}
	}

	ret := struct {
		Original string `json:"original"`
		Url      string `json:"url"`
		Title    string `json:"title"`
		State    string `json:"state"`
	}{Original, Url, Title, State}
	return c.RenderJson(ret)
}

func (c Admin) SaveProductDetail(id int64, content string) revel.Result {
	err := c.productApi().SaveProductDetail(id, content)
	if err != nil {
		return c.RenderJson(Error("保存信息出错,"+err.Error(), nil))
	}

	return c.RenderJson(Error("保存信息成功！", nil))
}

func (c Admin) UploadProductMaterial(id int64) revel.Result {
	count := 0
	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			to := ""
			p := entity.ProductParams{Type: entity.PTMaterial, Name: fileHeader.Filename, Value: to, ProductId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to = fmt.Sprintf("%d-%s", p.Id, fileHeader.Filename)
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			out, err := os.Create("data/products/m/" + to)
			gotang.AssertNoError(err, `os.Create`)

			in, err := fileHeader.Open()
			gotang.AssertNoError(err, `fileHeader.Open()`)

			io.Copy(out, in)

			out.Close()
			in.Close()
			count += 1
		}
	}
	if count == 0 {
		return c.RenderJson(Error("请选择要上传的文件", nil))
	}

	return c.RenderJson(Success("上传成功！", nil))
}

func (c Admin) deleteProductParams(id int64) revel.Result {
	if ret := c.checkErrorAsJsonResult(c.productApi().DeleteProductParams(id)); ret != nil {
		return ret
	}

	return c.RenderJson(Success("删除成功！", ""))
}

func (c Admin) DeleteSdImage(id int64) revel.Result {
	return c.deleteProductParams(id)
}

func (c Admin) DeleteImagePic(id int64) revel.Result {
	return c.deleteProductParams(id)
}

func (c Admin) DeleteMFile(id int64) revel.Result {
	return c.deleteProductParams(id)
	//TODO delete file?
}

func (c Admin) DoSaveProductSpec(productId int64, id int64, name string, value string) revel.Result {
	pp := entity.ProductParams{ProductId: productId, Id: id, Name: name, Value: value, Type: entity.PTSpec}
	if id == 0 { //new
		c.db.Insert(&pp)
	} else { //update
		c.db.Id(id).Update(&pp)
	}
	return c.RenderJson(Success("操作完成！", ""))
}

func (c Admin) DeleteSpec(id int64) revel.Result {
	return c.deleteProductParams(id)
}

func (c Admin) DeletePrice(id int64) revel.Result {
	var price entity.ProductPrices
	_, _ = c.db.Where("id=?", id).Get(&price)
	c.db.Delete(&price)
	return c.RenderJson(Success("", ""))
}

func (c Admin) ProductStockLogs(id int64) revel.Result {
	logs := c.productApi().FindAllProductStockLogs(id)
	return c.RenderJson(Success("", logs))
}

func (c Admin) AddProductStock(productId int64, stock int, message string) revel.Result {
	c.Validation.Required(stock != 0)
	if c.Validation.HasErrors() {
		return c.RenderJson(Error("请填入合法的入库数", nil))
	}
	newStock, err := c.productApi().AddProductStock(productId, stock, message)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("操作成功！", newStock))
}

func (c Admin) DoSaveProductPrice(productId int64, id int64, name string, price float64, start_quantity int, end_quantity int) revel.Result {
	if price <= 0 {
		return c.RenderJson(Error("请输入合法的价格(>=0)", "price"))
	}

	if start_quantity > 0 && end_quantity > 0 && end_quantity < start_quantity {
		return c.RenderJson(Error("起始价不应大于结束价", "start_quantity"))
	}

	var p entity.ProductPrices
	if id == 0 { //new
		p.ProductId = productId
		p.Name = name
		p.Price = price
		p.StartQuantity = start_quantity
		p.EndQuantity = end_quantity
		c.db.Insert(&p)
	} else {
		_, err := c.db.Where("id=?", id).Get(&p)
		if err != nil {
			return c.RenderJson(Error("操作失败！", err.Error()))
		}
		p.ProductId = productId
		p.Name = name
		p.Price = price
		p.StartQuantity = start_quantity
		p.EndQuantity = end_quantity
		c.db.Id(p.Id).Cols("name", "price", "start_quantity", "end_quantity").Update(&p)
	}
	return c.RenderJson(Success("操作成功！", ""))
}

/////////////////////////////////////////////////////

func (c Admin) Categories() revel.Result {
	c.setChannel("products/categories")
	return c.Render()
}

func (c Admin) CategoriesData() revel.Result {
	page := c.productApi().FindAllCategoriesForPage(c.pageSearcher())
	return c.renderDTJson(page)
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
		return c.RenderJson(Error("分类不存在", nil))
	}

	err := api.ToggleCategoryEnabled(&p)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(Success("激活成功！", nil))
		}
		return c.RenderJson(Success("禁用成功！", nil))
	}
}

//////////////////////////////////////////////////////////////////
// Providers

func (c Admin) Providers() revel.Result {
	c.setChannel("providers/providers")
	return c.Render()
}

func (c Admin) ProvidersData(filter_status string, filter_tags string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}

		if len(filter_tags) > 0 {
			session.And("tags=?", filter_tags)
		}
	})
	page := c.productApi().FindAllProvidersForPage(ps)
	return c.renderDTJson(page)
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
		return c.RenderJson(Error("制造商不存在", nil))
	}

	err := api.ToggleProviderEnabled(&p)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(Success("激活成功！", nil))
		}
		return c.RenderJson(Success("禁用成功！", nil))
	}
}

func (c Admin) DeleteProvider(id int64) revel.Result {
	_ = c.productApi().DeleteProvider(id)
	return c.RenderJson(Success("删除成功！", nil))
}

// 上传Logo
func (c Admin) UploadProviderImage(id int64, image *os.File) revel.Result {
	c.Validation.Required(image != nil)
	if c.Validation.HasErrors() {
		return c.RenderJson(Error("请选择图片", nil))
	}
	p, exists := c.productApi().GetProviderById(id)
	if !exists {
		return c.RenderJson(Error("操作失败，制造商不存在", nil))
	}
	to := filepath.Join(revel.Config.StringDefault("dir.data.providers", "data/providers"), fmt.Sprintf("%d.jpg", p.Id))

	err := utils.MakeAndSaveFromReader(image, to, "fit", 99, 44)
	if ret := c.checkUploadError(err, "保存上传图片报错！"); ret != nil {
		return ret
	}

	return c.RenderJson(Success("上传成功", nil))
}

// 上传Logo
func (c Admin) UploadProductLogo(id int64, image *os.File) revel.Result {
	c.Validation.Required(image != nil)
	if c.Validation.HasErrors() {
		return c.RenderJson(Error("请选择图片", nil))
	}
	p, exists := c.productApi().GetProductById(id)
	if !exists {
		return c.RenderJson(Error("操作失败，产品不存在", nil))
	}
	to := filepath.Join(revel.Config.StringDefault("dir.data.products/logo", "data/products/logo"), fmt.Sprintf("%d.jpg", p.Id))

	err := utils.MakeAndSaveFromReader(image, to, "thumbnail", 200, 200)
	if ret := c.checkUploadError(err, "保存上传图片报错！"); ret != nil {
		return ret
	}

	return c.RenderJson(Success("上传成功", nil))
}

func (c Admin) ProviderRecommends() revel.Result {
	c.setChannel("providers/recommends")
	return c.Render()
}

func (c Admin) ProductHots() revel.Result {
	c.setChannel("products/hots")
	return c.Render()
}

func (c Admin) checkUploadError(err error, msg string) revel.Result {
	if err != nil {
		revel.WARN.Printf("上传头像操作失败，%s， msg：%s", msg, err.Error())
		return c.RenderJson(Error("操作失败，"+msg+", "+err.Error(), nil))
	}
	return nil
}

func (c Admin) AdImages() revel.Result {
	c.setChannel("system/adimages")
	return c.Render()
}

func (c Admin) UploadAdImage() revel.Result {
	//698, 191
	var (
		dir   = "data/adimages/"
		t     = entity.ATAd
		count = 0
	)

	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			to := utils.Uuid() + ".jpg"
			p := entity.AppParams{Type: t, Name: fileHeader.Filename, Value: to}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, "fit", 698, 191)
			gotang.AssertNoError(err, "生成图片出错！")

			count += 1
		}
	}

	if count == 0 {
		return c.RenderJson(Error("请选择要上传的图片", nil))
	}

	return c.RenderJson(Success("上传成功！", nil))
}

func (c Admin) DeleteAdImage(id int64) revel.Result {
	c.appApi().DeleteAdImage(id)
	return c.RenderJson(Success("", ""))
}

func (c Admin) SetFirstAdImageUrl(id int64) revel.Result {
	_ = c.appApi().SetFirstAdImage(id)
	return c.RenderJson(Success("", ""))
}

func (c Admin) HotKeywords() revel.Result {
	c.setChannel("system/hotkeywords")
	return c.Render()
}

func (c Admin) DeleteHotKeyword(id int64) revel.Result {
	c.appApi().DeleteHotKeyword(id)
	return c.RenderJson(Success("", ""))
}

func (c Admin) SetFirstHotKeyword(id int64) revel.Result {
	_ = c.appApi().SetFirstHotKeyword(id)
	return c.RenderJson(Success("", ""))
}

func (c Admin) DoSaveHotKeyword(id int64, value string) revel.Result {
	pp := entity.AppParams{Id: id, Name: "", Value: value, Type: entity.ATHk}
	if id == 0 { //new
		c.db.Insert(&pp)
	} else { //update
		c.db.Id(id).Update(&pp)
	}
	return c.RenderJson(Success("操作完成！", ""))
}

func (c Admin) Slogan() revel.Result {
	c.setChannel("system/slogan")
	p, _ := c.appApi().GetSlogan()
	return c.Render(p)
}

func (c Admin) SaveSlogan(p entity.AppParams) revel.Result {
	c.appApi().SaveSlogan(p)

	return c.Redirect(Admin.Slogan)
}

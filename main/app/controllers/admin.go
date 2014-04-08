package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/itang/gotang"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
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
	userTotal := c.userApi().Total()
	orderTotal := c.orderApi().TotalNewOrders()

	c.setChannel("/")
	return c.Render(userTotal, orderTotal)
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

func (c Admin) ToggleCommentEnabled(id int64) revel.Result {
	comment, ok := c.userApi().GetCommentById(id)
	if !ok {
		return c.RenderJson(Error("评论不存在", nil))
	}

	err := c.userApi().ToggleCommentEnabled(&comment)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if comment.Enabled {
			return c.RenderJson(Success("审核评论通过！", nil))
		}
		return c.RenderJson(Success("审核评论不通过！", nil))
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
func (c Admin) ProductsData(filter_status string, filter_tag string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
		if len(filter_tag) > 0 {
			session.And("tags like ?", "%"+filter_tag+"%")
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
		splits    = ""
	)

	if id == 0 { // new
		p.MinNumberOfOrders = 1
	} else { //edit
		p, _ = c.productApi().GetProductById(id)
		detail, _ = c.productApi().GetProductDetail(p.Id)
		stockLogs = c.productApi().FindAllProductStockLogs(p.Id)
		splits = c.productApi().GetProductPricesSplits(p.Id)
	}
	revel.INFO.Println("splits", splits)
	return c.Render(p, detail, stockLogs, splits)
}

func (c Admin) DoNewProduct(p entity.Product) revel.Result {
	c.Validation.Required(p.Name).Message("请填写名称").Key("name")
	c.Validation.Required(p.MinNumberOfOrders >= 1).Message("起订最小数量应该大于0").Key("min_number_of_orders")

	if ret := c.doValidate(routes.Admin.NewProduct(p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveProduct(p)
	if err != nil {
		c.Flash.Error("保存产品失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存产品成功！")
	}

	return c.Redirect(routes.Admin.NewProduct(id))
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

	return c.RenderJson(Success("保存信息成功！", nil))
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

func (c Admin) DoSaveSplitProductPrice(productId int64, start_quantitys string) revel.Result {
	err := c.productApi().SplitProductPrices(productId, start_quantitys)
	if err != nil {
		return c.RenderJson(Error(err.Error(), "start_quantitys"))
	}
	return c.RenderJson(Success("操作成功！", ""))
}

func (c Admin) DoSaveProductPrice(productId int64, id int64, price float64) revel.Result {
	if price <= 0 {
		return c.RenderJson(Error("请输入合法的价格(>=0)", "price"))
	}

	var p entity.ProductPrices
	_, err := c.db.Where("id=?", id).Get(&p)
	if err != nil {
		return c.RenderJson(Error("操作失败！", err.Error()))
	}
	p.Price = price
	_, err = c.db.Id(p.Id).Cols("price").Update(&p)

	//更新冗余的价格
	err = c.productApi().UpdateProductPrice(productId)
	gotang.AssertNoError(err, "")

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

	if ret := c.doValidate(routes.Admin.NewCategory(p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveCategory(p)
	if err != nil {
		c.Flash.Error("保存分类失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存分类成功！")
	}

	return c.Redirect(routes.Admin.NewCategory(id))
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

	if ret := c.doValidate(routes.Admin.NewProvider(p.Id)); ret != nil {
		return ret
	}

	id, err := c.productApi().SaveProvider(p)
	if err != nil {
		c.Flash.Error("保存制造商失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存制造商成功！")
	}

	return c.Redirect(routes.Admin.NewProvider(id))
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

func (c Admin) SetAdImageLink(id int64, link string) revel.Result {
	_ = c.appApi().SetAdImageLink(id, link)
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

func (c Admin) SelfDelivery() revel.Result {
	c.setChannel("system/self_delivery")
	return c.Render()
}

func (c Admin) Payments() revel.Result {
	c.setChannel("system/payments")
	return c.Render()
}

func (c Admin) Orders() revel.Result {
	osJSON := utils.ToJSON(entity.OSMap)
	pmJSON := utils.ToJSON(entity.PMMap)
	spJSON := utils.ToJSON(entity.SPMap)

	c.setChannel("orders/index")
	return c.Render(osJSON, pmJSON, spJSON)
}

func (c Admin) OrdersData(filter_status int) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		if filter_status != 0 {
			session.And("status=?", filter_status)
		}
	})
	orders := c.orderApi().FindSubmitedOrdersForPage(ps)
	return c.renderDTJson(orders)
}

func (c Admin) ShowOrder(userId int64, code int64) revel.Result {
	order, exists := c.orderApi().GetOrder(userId, code)
	if !exists {
		return c.NotFound("订单不存在!")
	}
	orderBy := c.userApi().GetUserDesc(order.UserId)
	return c.Render(order, orderBy)
}

func (c Admin) ToggleOrderLock(id int64) revel.Result {
	err := c.orderApi().ToggleOrderLock(id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c Admin) ChangeOrderPayed(id int64) revel.Result {
	err := c.orderApi().ChangeOrderPayed(id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c Admin) ChangeOrderVerify(id int64) revel.Result {
	err := c.orderApi().ChangeOrderVerify(id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c Admin) ChangeOrderShiped(id int64) revel.Result {
	err := c.orderApi().ChangeOrderShiped(id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c Admin) ProductComments() revel.Result {
	c.setChannel("products/comments")
	return c.Render()
}

func (c Admin) ProductCommentsData(filter_status string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
		session.And("target_type=?", entity.CT_PRODUCT)
	})
	page := c.userApi().CommentsForPage(ps)
	return c.renderDTJson(page)
}

func (c Admin) Prices() revel.Result {
	c.setChannel("prices/index")
	return c.Render()
}

func (c Admin) PricesData(filter_status string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		if filter_status == "1" {
			session.And("replies = 0")
		} else if filter_status == "2" {
			session.And("replies > 0")
		}
	})

	page := c.orderApi().FindAllInquiresForPage(ps)
	return c.renderDTJson(page)
}

func (c Admin) NewInquiryReply(id int64) revel.Result {
	in, exists := c.orderApi().GetInquiryById(id)
	if !exists {
		return c.NotFound("此询价不存在！")
	}

	replies := c.orderApi().GetInquiryReplies(id)
	return c.Render(in, replies)
}

func (c Admin) DoNewInquiryReply(reply entity.InquiryReply) revel.Result {
	reply.UserId = c.forceSessionUserId()
	err := c.orderApi().SaveInquiryReply(reply)

	if err != nil {
		c.FlashParams()
		c.Flash.Error("回复出错，请重试！")
		return c.Redirect(routes.Admin.NewInquiryReply(reply.InquiryId))
	}
	c.Flash.Success("回复成功！")
	return c.Redirect(routes.Admin.NewInquiryReply(reply.InquiryId))
}

func (c Admin) DeleteInquiryReply(id int64) revel.Result {
	err := c.orderApi().DeleteInquiryReply(id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("操作完成", ""))
}

func (c Admin) Site() revel.Result {
	ps := c.appConfigApi().FindConfigsBySection("site.basic")
	c.setChannel("system/site")
	return c.Render(ps)
}

func (c Admin) SaveSiteBasic(p []entity.StringKV) revel.Result {
	c.Flash.Success("保存成功！")
	for _, v := range p {
		c.appConfigApi().SaveOrUpdateConfig(v.Key, v.Value, "")
	}
	return c.Redirect(Admin.Site)
}

func (c Admin) SiteComment() revel.Result {
	ps := c.appConfigApi().FindConfigsBySection("site.comment")
	c.setChannel("system/system_comment")
	return c.Render(ps)
}

func (c Admin) SaveSiteComment(p []entity.StringKV) revel.Result {
	c.Flash.Success("保存成功！")
	for _, v := range p {
		c.appConfigApi().SaveOrUpdateConfig(v.Key, v.Value, "")
	}
	return c.Redirect(Admin.SiteComment)
}

func (c Admin) Contact() revel.Result {
	ps := c.appConfigApi().FindConfigsBySection("site.contact")
	c.setChannel("system/contact")
	return c.Render(ps)
}

func (c Admin) SaveSiteContact(p []entity.StringKV) revel.Result {

	c.Flash.Success("保存成功！")
	for _, v := range p {
		c.appConfigApi().SaveOrUpdateConfig(v.Key, v.Value, "")
	}
	return c.Redirect(Admin.Contact)
}

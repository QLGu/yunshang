package controllers

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/kr/pretty"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 用户相关Actions
type User struct {
	ShouldLoginedController
}

// 用户主页
func (c User) Index() revel.Result {
	currUser, _ := c.currUser()

	userLevel, _ := c.userApi().GetUserLevel(&currUser)
	revel.INFO.Printf("%v", userLevel)

	collects := c.userApi().TotalUserCollects(currUser.Id)
	carts := c.orderApi().TotalUserCarts(currUser.Id)

	topays := c.orderApi().TotalUserOrdersByStatus(currUser.Id, entity.OS_SUBMIT)

	receives := c.orderApi().TotalUserOrdersByStatus(currUser.Id, entity.OS_SHIP)

	c.setChannel("user/index")
	return c.Render(currUser, userLevel, collects, carts, topays, receives)
}

// 到用户信息
func (c User) UserInfo() revel.Result {
	c.setChannel("user/userinfo")

	c.FlashParams()

	user := c.forceCurrUser()
	userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)
	return c.Render(user, userDetail)
}

// 保存用户信息
func (c User) DoSaveUserInfo(user entity.User, userDetail entity.UserDetail) revel.Result {
	revel.INFO.Printf("user: %v", user)
	revel.INFO.Printf("userDetail: %v", userDetail)

	currUser := c.forceCurrUser()

	c.Validation.Match(user.MobilePhone, regexp.MustCompile(`^(1(([35][0-9])|(47)|[8][01236789]))\d{8}$`)).Message("请填入正确的手机号码")

	if userDetail.CompanyWebsite != "" {
		c.Validation.Match(userDetail.CompanyWebsite, regexp.MustCompile(`http://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`)).Message("请填入正确的网址")
	}

	if ret := c.doValidate(User.UserInfo); ret != nil {
		return ret
	}

	if len(currUser.LoginName) == 0 {
		c.Validation.Required(user.LoginName).Message("请输入登录名")
		c.Validation.MinSize(user.LoginName, 4).Message("请输入至少4位登录名")
		c.Validation.MaxSize(user.LoginName, 100).Message("输入过多位数的登录名")

		ok := true
		for _, v := range []rune(user.LoginName) {
			if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') || v == '_') {
				ok = false
				break
			}
		}
		c.Validation.Required(ok).Message("输入的登录名不符合要求")

		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}

		_, exists := c.userApi().CheckUserByLoginName(user.LoginName)
		c.Validation.Required(!exists).Message("该用户名已经注册").Key("user.LoginName")
	}

	if len(currUser.Email) == 0 {
		revel.INFO.Println("len(currUser.Email) == 0 " + user.Email)
		c.Validation.Required(user.Email).Message("请输入邮箱")
		c.Validation.Email(user.Email).Message("请输入合法邮箱")
		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}
		_, exists := c.userApi().CheckUserByEmail(user.Email)
		c.Validation.Required(!exists).Message("该邮箱已经注册").Key("user.Email")
		if ret := c.doValidate(User.UserInfo); ret != nil {
			return ret
		}
	}

	err := c.userApi().UpdateUserInfo(&currUser, user, userDetail)
	if err != nil {
		c.Flash.Error("保存用户信息失败" + err.Error())
	} else {
		c.Flash.Success("保存会员信息成功！")
	}
	return c.Redirect(User.UserInfo)
}

// 到修改密码
func (c User) ChangePassword() revel.Result {
	hasPassword := len(c.forceCurrUser().CryptedPassword) != 0

	c.setChannel("user/userinfo/cw")

	return c.Render(hasPassword)
}

// 修改密码处理
func (c User) DoChangePassword(oldPassword, password, confirmPassword string) revel.Result {
	c.Validation.Required(oldPassword).Message("请输入原始密码")
	c.Validation.Required(password).Message("请输入新密码")
	c.Validation.MinSize(password, 6).Message("请输入6位新密码")
	c.Validation.MaxSize(password, 50).Message("输入新密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认新密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的新密码不匹配").Key("confirmPassword")
	c.Validation.Required(password != oldPassword).Message("输入了与原始密码相同的新密码").Key("password")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	user := c.forceCurrUser()
	c.Validation.Required(c.userApi().VerifyPassword(user.CryptedPassword, oldPassword)).Message("你的原始密码输入有误").Key("oldPassword")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	if err := c.userApi().DoChangePassword(&user, password); err != nil {
		c.Flash.Error("修改密码失败：" + err.Error())
	} else {
		c.Flash.Success("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
	}

	return c.Redirect(User.ChangePassword)
}

// 到设置密码
func (c User) SetPassword(password, confirmPassword string) revel.Result {
	c.Validation.Required(password).Message("请输入密码")
	c.Validation.MinSize(password, 6).Message("请输入6位密码")
	c.Validation.MaxSize(password, 50).Message("输入密码位数太长了")
	c.Validation.Required(confirmPassword).Message("请输入确认密码")
	c.Validation.Required(password == confirmPassword).Message("两次输入的密码不匹配").Key("confirmPassword")
	if ret := c.doValidate(User.ChangePassword); ret != nil {
		return ret
	}

	user := c.forceCurrUser()
	if err := c.userApi().DoChangePassword(&user, password); err != nil {
		c.Flash.Error("修改密码失败：" + err.Error())
	} else {
		c.Flash.Success("修改密码成功，你的新密码是：" + password[0:3] + strings.Repeat("*", len(password)-5) + password[len(password)-2:])
	}

	return c.Redirect(User.ChangePassword)
}

// 用户级别显示
func (c User) UserLevel() revel.Result {
	currUser, _ := c.currUser()
	userLevel, _ := c.userApi().GetUserLevel(&currUser)

	userLevels := c.userApi().FindUserLevels()
	userScores := currUser.Scores

	c.setChannel("points/level")
	return c.Render(userLevels, userLevel, userScores)
}

// 积分规则显示
func (c User) ScoresRules() revel.Result {
	currUser, _ := c.currUser()
	userLevel, _ := c.userApi().GetUserLevel(&currUser)

	userLevels := c.userApi().FindUserLevels()
	userScores := currUser.Scores

	c.setChannel("points/rules")
	return c.Render(userLevels, userLevel, userScores)
}

// 用户的订单
func (c User) Orders() revel.Result {
	orders := c.orderApi().FindSubmitedOrdersByUser(c.forceSessionUserId())

	c.setChannel("order/orders")
	return c.Render(orders)
}

// 显示用户头像
// param file： 头像图片标识： {{ucode}}.jpg
func (c User) Image(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.iamges", "data/images")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No Image Found！")
	}

	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

// 上传用户头像
func (c User) UploadImage(image *os.File) revel.Result {
	c.Validation.Required(image != nil)
	if ret := c.doValidate(User.UploadImage); ret != nil {
		return c.RenderJson(Error("请选择图片", nil))
	}
	ucode, ok := c.Session["ucode"]
	gotang.Assert(ok, "ucode")

	from := image.Name()
	to := filepath.Join(revel.Config.StringDefault("dir.data.iamges", "data/images"), ucode+".jpg")

	err := utils.MakeAndSaveThumbnail(from, to, 460, 460)
	if ret := c.checkUploadError(err, "保存上传图片报错！"); ret != nil {
		return ret
	}

	return c.RenderJson(Success("上传成功", nil))
}

func (c User) checkUploadError(err error, msg string) revel.Result {
	if err != nil {
		revel.WARN.Printf("上传头像操作失败，%s， msg：%s", msg, err.Error())
		return c.RenderJson(Error("操作失败，"+msg+", "+err.Error(), nil))
	}
	return nil
}

// 用户收货地址列表
func (c User) DeliveryAddresses() revel.Result {
	das := c.userApi().FindUserDeliveryAddresses(c.forceSessionUserId())

	c.setChannel("userinfo/das")
	return c.Render(das)
}

func (c User) DeliveryAddressesData() revel.Result {
	das := c.userApi().FindUserDeliveryAddresses(c.forceSessionUserId())
	return c.RenderJson(Success("", das))
}

func (c User) DasForSelect() revel.Result {
	das := c.userApi().FindUserDeliveryAddresses(c.forceSessionUserId())
	return c.Render(das)
}

func (c User) NewDeliveryAddress(id int64) revel.Result {
	revel.INFO.Println("NewDeliveryAddress, id", id)
	user := c.forceCurrUser()

	var da entity.DeliveryAddress
	if id == 0 { // new
		revel.INFO.Println("total", c.userApi().GetUserDeliveryAddressTotal(user.Id))
		if c.userApi().GetUserDeliveryAddressTotal(user.Id) == 0 {
			userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)
			da = entity.DeliveryAddress{
				Consignee:   user.RealName,
				MobilePhone: user.MobilePhone,
				Province:    userDetail.CompanyProvince,
				City:        userDetail.CompanyCity,
				Area:        userDetail.CompanyArea,
				Address:     userDetail.CompanyAddress,
				Email:       user.Email,
			}
		}
	} else { //edit
		da, _ = c.userApi().GetUserDeliveryAddress(user.Id, id)
	}
	revel.INFO.Printf("%v", da)
	return c.Render(da)
}

func (c User) DoNewDeliveryAddress(da entity.DeliveryAddress) revel.Result {
	revel.INFO.Println("da.ismain", da.IsMain)

	c.Validation.Required(da.Name).Message("请输入收货地址命名")
	c.Validation.Required(da.Consignee).Message("请输入收货人")
	c.Validation.Required(len(da.MobilePhone) > 0 || len(da.FixedPhone) > 0).Message("请输入电话号码").Key("da.MobilePhone")
	c.Validation.Required(da.Province).Message("请输入省")
	c.Validation.Required(da.City).Message("请输入城市")
	c.Validation.Required(da.Area).Message("请输入地区")
	if len(da.Email) != 0 {
		c.Validation.Email(da.Email).Message("请输入合法的邮箱")
	}

	if len(da.MobilePhone) != 0 {
		c.Validation.Match(da.MobilePhone, regexp.MustCompile(`^(1(([35][0-9])|(47)|[8][01236789]))\d{8}$`)).Message("请填入正确的手机号码")
	}

	if len(da.FixedPhone) != 0 {
		c.Validation.Match(da.FixedPhone, regexp.MustCompile(`^0\d{2,3}(\-)?\d{7,8}$`)).Message("请填入正确的固定电话")
	}

	if ret := c.doValidate(routes.User.NewDeliveryAddress(da.Id)); ret != nil {
		return ret
	}

	da.UserId = c.forceSessionUserId()
	daId, err := c.userApi().SaveUserDeliveryAddress(da)
	if err != nil {
		c.Flash.Error("保存收货地址失败，请重试！")
	} else {
		c.Flash.Success("保存收货地址成功！")
	}
	revel.INFO.Printf("daid:%v", daId)
	return c.Redirect(routes.User.NewDeliveryAddress(daId))
}

func (c User) DeleteDeliveryAddress(id int64) revel.Result {

	_ = c.userApi().DeleteDeliveryAddress(c.forceSessionUserId(), id)

	return c.RenderJson(Success("ok", nil))
}

func (c User) Cart() revel.Result {
	carts := c.orderApi().FindUserCarts(c.forceSessionUserId())
	c.setChannel("order/cart")
	return c.Render(carts)
}

func (c User) CartData() revel.Result {
	carts := c.orderApi().FindUserCarts(c.forceSessionUserId())
	return c.RenderJson(Success("", carts))
}

func (c User) CartProductPrices() revel.Result {
	ps := c.orderApi().FindUserCartProductPrices(c.forceSessionUserId())
	return c.RenderJson(Success("", ps))
}

func (c User) OrderProducts(code int64) revel.Result {
	ps := c.orderApi().FindOrderProducts(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) CartProductPrefPrice(productId int64, quantity int) revel.Result {
	ps := c.productApi().GetProductPrefPrice(productId, quantity)
	return c.RenderJson(Success("", ps))
}

func (c User) Prices() revel.Result {
	c.setChannel("order/prices")
	return c.Render()
}

func (c User) Comments() revel.Result {
	c.setChannel("userinfo/comments")
	return c.Render()
}

//////////////////////////////////////////////////////////////////////////////////////
// collects
func (c User) Collects() revel.Result {
	c.setChannel("userinfo/collects")
	return c.Render()
}

func (c User) CollectsData() revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		session.And("user_id=?", c.forceSessionUserId())
	})
	page := c.userApi().FindAllProductCollectsForPage(ps)
	// 加上当前价格
	collects := page.Data.([]entity.ProductCollect)

	type ppw struct {
		entity.ProductCollect
		CurrentPrice float64 `json:"current_price"`
		Name         string  `json:"name"`
	}

	ppws := make([]ppw, len(collects))
	for i, v := range collects {
		p, _ := c.productApi().GetProductById(v.ProductId)
		ppws[i] = ppw{v, p.Price, p.Name}
	}
	page.Data = ppws

	return c.renderDTJson(page)
}

func (c User) CollectProduct(id int64) revel.Result {
	err := c.userApi().CollectProduct(c.forceSessionUserId(), id)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("收藏产品成功！", nil))
}

func (c User) DeleteProductCollect(id int64) revel.Result {
	userId := c.forceSessionUserId()
	var p entity.ProductCollect
	ok, err := c.db.Where("id=? and user_id=?", id, userId).Get(&p)
	if !ok || err != nil {
		return c.RenderJson(Error("此收藏不存在！", ""))
	}
	_, _ = c.db.Delete(&p)
	return c.RenderJson(Success("删除收藏成功！", nil))
}

func (c User) AddToCart(productId int64) revel.Result {
	err := c.orderApi().AddProductToCart(c.forceSessionUserId(), productId, 1)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c User) AddToCartResult(productId int64) revel.Result {
	return c.Render()
}

func (c User) DeleteCartProduct(id int64) revel.Result {
	err := c.orderApi().DeleteCartProduct(id, c.forceSessionUserId())
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c User) CleanCart() revel.Result {
	err := c.orderApi().CleanCart(c.forceSessionUserId())
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c User) MoveCartsToCollects() revel.Result {
	err := c.orderApi().MoveCartsToCollects(c.forceSessionUserId())
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", ""))
}

func (c User) IncCartProductQuantity(id int64, quantity int) revel.Result {
	p, err := c.orderApi().IncCartProductQuantity(id, c.forceSessionUserId(), quantity)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}

	return c.RenderJson(Success("", p))
}

func (c User) DoNewOrder(ps []entity.ParamsForNewOrder) revel.Result {
	revel.INFO.Printf("%v", ps)
	filterPs := make([]entity.ParamsForNewOrder, 0)
	for _, p := range ps {
		if p.CartId != 0 {
			filterPs = append(filterPs, p)
		}
	}
	c.Validation.Required(len(filterPs) > 0).Message("请选择商品！")
	if ret := c.doValidate(User.Cart); ret != nil {
		return ret
	}

	order, err := c.orderApi().SaveTempOrder(c.forceSessionUserId(), filterPs)
	c.Validation.Required(err == nil).Message("保存订单出错， 请重试！")
	if ret := c.doValidate(User.Cart); ret != nil {
		c.setRollbackOnly()
		revel.ERROR.Printf("%v", err)
		return ret
	}

	return c.Redirect(routes.User.ConfirmOrder(order.Code))
}

func (c User) ConfirmOrder(code int64) revel.Result {
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}
	if order.IsSubmited() {
		c.Flash.Error("订单不存在！")
		return c.Redirect(User.Cart)
	}

	shippings := c.orderApi().FindShippings(order.Amount)
	payments := c.orderApi().FindAPayments()
	c.setChannel("order/confirm")
	return c.Render(order, shippings, payments)
}

func (c User) OrderItems(code int64) revel.Result {
	ps := c.orderApi().GetOrderItems(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) OrderData(code int64) revel.Result {
	ps, _ := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	return c.RenderJson(Success("", ps))
}

func (c User) DoSubmitOrder(o entity.Order) revel.Result {
	revel.INFO.Printf("%# v", pretty.Formatter(o))
	err := c.orderApi().SubmitOrder(c.forceSessionUserId(), o)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(routes.User.ConfirmOrder(o.Code))
	}

	c.Flash.Success("提交订单成功!")
	return c.Redirect(routes.User.PayOrder(o.Code))
}

func (c User) PayOrder(code int64) revel.Result {
	c.setChannel("order/orders/pay")
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}
	return c.Render(order)
}

func (c User) CancelOrder(code int64) revel.Result {
	err := c.orderApi().CancelOrderByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) ViewOrder(code int64) revel.Result {
	c.setChannel("order/orders/view")
	order, exists := c.orderApi().GetOrder(c.forceSessionUserId(), code)
	if !exists {
		return c.NotFound("订单不存在!")
	}

	orderBy := c.userApi().GetUserDesc(order.UserId)

	return c.Render(order, orderBy)
}

func (c User) DeleteOrder(code int64) revel.Result {
	err := c.orderApi().DeleteOrderByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) PayOrderByUserComment(code int64, comment string) revel.Result {
	err := c.orderApi().PayOrderByUserComment(c.forceSessionUserId(), code, comment)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ""))
}

func (c User) OrderLogsData(code int64) revel.Result {
	ps, err := c.orderApi().FindAllOrderLogsByUser(c.forceSessionUserId(), code)
	if ret := c.checkErrorAsJsonResult(err); ret != nil {
		return ret
	}
	return c.RenderJson(Success("", ps))
}

func (c User) Invoices() revel.Result {
	ins := c.userApi().FindUserInvoices(c.forceSessionUserId())

	c.setChannel("userinfo/ins")
	return c.Render(ins)
}

func (c User) NewInvoice(id int64) revel.Result {
	user := c.forceCurrUser()

	var in entity.Invoice
	if id == 0 { // new
		userDetail, _ := c.userApi().GetUserDetailByUserId(user.Id)
		in = entity.Invoice{
			CompanyName:    userDetail.CompanyName,
			CompanyAddress: userDetail.CompanyFullAddress(),
			CompanyPhone:   userDetail.CompanyPhone,
			DaZipCode:      userDetail.CompanyZipCode,
			Type:           entity.IN_COMMON,
		}
	} else { //edit
		in, _ = c.userApi().GetUserInvoice(user.Id, id)
	}
	return c.Render(in)
}

func (c User) DoNewInvoice(in entity.Invoice) revel.Result {
	c.Validation.Required(in.Type != 0).Message("请选择发票类型！").Key("in.Type")
	if ret := c.doValidate(routes.User.NewInvoice(in.Id)); ret != nil {
		return ret
	}

	in.UserId = c.forceSessionUserId()
	id, err := c.userApi().SaveUserInvoice(in)
	if err != nil {
		c.Flash.Error("保存发票信息，请重试！")
	} else {
		c.Flash.Success("保存发票信息成功！")
	}
	return c.Redirect(routes.User.NewInvoice(id))
}

func (c User) DeleteInvoice(id int64) revel.Result {
	_ = c.userApi().DeleteInvoice(c.forceSessionUserId(), id)

	return c.RenderJson(Success("ok", nil))
}

func (c User) InsForSelect() revel.Result {
	ins := c.userApi().FindUserInvoices(c.forceSessionUserId())
	return c.Render(ins)
}

func (c User) OrderDaForView(order entity.Order) revel.Result {
	da, ok := c.orderApi().GetDaForView(order.UserId, order.DaId)
	return c.Render(da, ok)
}

func (c User) OrderInForView(order entity.Order) revel.Result {
	in, ok := c.orderApi().GetInForView(order.UserId, order.InvoiceId)
	return c.Render(in, ok)
}

func (c User) OrderShippingForView(order entity.Order) revel.Result {
	shipping, ok := c.orderApi().GetShippingForView(order.UserId, order.ShippingId)
	return c.Render(shipping, ok)
}

func (c User) OrderPaymentForView(order entity.Order) revel.Result {
	payment, ok := c.orderApi().GetPaymentForView(order.UserId, order.PaymentId)
	return c.Render(payment, ok)
}

func (c User) ReceiptOrder(code int64) revel.Result {
	_ = c.orderApi().ReceiptOrder(c.forceSessionUserId(), code)

	return c.RenderJson(Success("ok", nil))
}

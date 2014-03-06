package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"

	"github.com/itang/gotang"
	reveltang "github.com/itang/reveltang/controllers"
	"github.com/itang/yunshang/main/app/models"
	"github.com/itang/yunshang/main/app/models/entity"
)

// Rest响应的数据结构
type RestResposne struct {
	Ok      bool        `json:"ok"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// DataTables server-side响应数据结构
type dataTableData struct {
	SEcho                int         `json:"sEcho"`
	ITotalRecords        int64       `json:"iTotalRecords"`
	ITotalDisplayRecords int64       `json:"iTotalDisplayRecords"`
	AaData               interface{} `json:"aaData,omitempty"`
}

// 构建dataTableData
func DataTableData(echo string, total int64, totalDisplay int64, data interface{}) dataTableData {
	ei, err := strconv.Atoi(echo)
	if err != nil {
		ei = 0
	}
	return dataTableData{SEcho: ei, ITotalRecords: total, ITotalDisplayRecords: totalDisplay, AaData: data}
}

// 应用控制器
type AppController struct {
	*revel.Controller
	XOrmTnController
	reveltang.XRuntimeableController
}

// 初始化逻辑
func (c AppController) init() revel.Result {
	c.RenderArgs["_host"] = revel.Config.StringDefault("web.host", "localhost:9000")
	c.setChannel("")

	return nil
}

// 成功的Response
func (c AppController) successResposne(message string, data interface{}) RestResposne {
	return RestResposne{Ok: true, Code: 0, Message: message, Data: data}
}

// 失败的Response
func (c AppController) errorResposne(message string, data interface{}) RestResposne {
	return RestResposne{Ok: false, Code: 1, Message: message, Data: data}
}

// 是否登录?
func (c AppController) isLogined() bool {
	_, ok := c.Session["uid"]
	return ok
}

// 获取当前用户
func (c AppController) currUser() (user entity.User, ok bool) {
	uidStr, ok := c.Session["uid"]
	if !ok {
		return user, false
	}

	id, err := strconv.Atoi(uidStr)
	if err != nil {
		return user, false
	}
	return c.userService().GetUserById(int64(id))
}

// 获取当前Session用户
func (c AppController) currSessionUser() (user models.SessionUser, ok bool) {
	uidStr, ok := c.Session["uid"]
	if !ok {
		return user, false
	}
	id, err := strconv.Atoi(uidStr)
	if err != nil {
		return user, false
	}

	user.Id = int64(id)
	user.Email, _ = c.Session["email"]
	user.From, _ = c.Session["from"]
	user.LoginName, _ = c.Session["login"]
	return
}

// 设置用户会话信息
func (c AppController) setLoginSession(sessionUser models.SessionUser) {
	c.Session["login"] = sessionUser.LoginName
	c.Session["uid"] = fmt.Sprintf("%v", sessionUser.Id)
	c.Session["screen_name"] = sessionUser.DisplayName()
	c.Session["email"] = sessionUser.Email
	c.Session["from"] = sessionUser.From
}

// 清理用户会话信息
func (c AppController) clearLoginSession() {
	delete(c.Session, "uid")
	delete(c.Session, "login")
	delete(c.Session, "screen_name")
	delete(c.Session, "email")
	delete(c.Session, "from")
}

// 强制获取当前登录用户，如果不存在则Panic
func (c AppController) forceCurrUser() (user entity.User) {
	user, ok := c.currUser()
	gotang.Assert(ok, "用户未登录！")
	return
}

// 执行验证逻辑
func (c AppController) doValidate(redirectTarget interface{}) revel.Result {
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(redirectTarget)
	}
	return nil
}

// 获取客户端IP
func (c AppController) getRemoteIp() string {
	ips, ok := c.Request.Header["X-Real-IP"]
	if !ok {
		return strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	return ips[0]
}

// 输出DataTable 分页数据
func (c AppController) renderDataTableJson(page models.PageData) revel.Result {
	sEcho := c.Params.Get("sEcho")
	return c.RenderJson(DataTableData(sEcho, page.Total, page.Total, page.Data))
}

// 构造分页查询器
func (c AppController) pageSearcher() *models.PageSearcher {
	var (
		start     int
		limit     int
		search    string
		sortColNo string
		sortField string
		sortDir   string
	)

	c.Params.Bind(&start, "iDisplayStart")
	c.Params.Bind(&limit, "iDisplayLength")
	c.Params.Bind(&search, "sSearch")
	c.Params.Bind(&sortColNo, "iSortCol_0")
	c.Params.Bind(&sortField, "mDataProp_"+sortColNo)
	c.Params.Bind(&sortDir, "sSortDir_0")
	if limit == 0 {
		limit = 10
	}
	return &models.PageSearcher{
		Start: start, Limit: limit,
		SortField: sortField, SortDir: sortDir,
		Search: search, Session: c.XOrmSession}
}

// 构造分页查询器，通过附加的Session回调
func (c AppController) pageSearcherWithCalls(calls ...models.PageSearcherCall) *models.PageSearcher {
	ps := c.pageSearcher()
	ps.OtherCalls = calls
	return ps
}

// 用户服务
func (c AppController) userService() models.UserService {
	gotang.Assert(c.XOrmSession != nil, "c.XOrmSession should no be nil")
	return models.DefaultUserService(c.XOrmSession)
}

// 设置当前页Channel
func (c AppController) setChannel(ch string) {
	c.RenderArgs["channel"] = ch
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ShouldLoginedController struct {
	AppController
}

// 检查用户是否登录，如果未登录，则转入主页
func (c ShouldLoginedController) checkUser() revel.Result {
	if !c.isLogined() {
		return c.Redirect(App.Index)
	}
	return nil
}

type AdminController struct {
	ShouldLoginedController
}

// 检查用户是否为管理员， 如果不是，则转入主页
func (c AdminController) checkAdminUser() revel.Result {
	user, _ := c.Session["screen_name"]
	//TODO 使用角色
	if !c.isLogined() || user != "admin" {
		return c.Redirect(App.Index)
	}
	return nil
}

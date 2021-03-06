package controllers

import (
	"fmt"
	"reflect"
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

// 成功的Response
func Success(message string, data interface{}) RestResposne {
	return newRestResponse(true, 0, message, data)
}

// 失败的Response
func Error(message string, data interface{}) RestResposne {
	return newRestResponse(false, 1, message, data)
}

func isEmptySlice(data interface{}) bool {
	return data != nil && (reflect.TypeOf(data).Kind() == reflect.Slice && reflect.ValueOf(data).IsNil())
}

func newRestResponse(ok bool, code int, message string, data interface{}) RestResposne {
	if isEmptySlice(data) {
		data = []string{}
	}
	return RestResposne{Ok: ok, Code: code, Message: message, Data: data}
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
	__userApi      *models.UserService
	__productApi   *models.ProductService
	__appApi       *models.AppService
	__orderApi     *models.OrderService
	__newsApi      *models.NewsService
	__appConfigApi *models.AppConfigService
}

// 初始化逻辑
func (c AppController) init() revel.Result {
	c.RenderArgs["_host"] = models.CacheSystem.GetConfig("site.basic.host")
	c.setChannel("")

	return nil
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
	return c.userApi().GetUserById(int64(id))
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
	user.Roles, _ = c.Session["roles"]
	return
}

func (c AppController) forceSessionUserId() int64 {
	uidStr, ok := c.Session["uid"]
	gotang.Assert(ok, "获取当前登录用户失败！")

	userId, err := strconv.Atoi(uidStr)
	gotang.AssertNoError(err, "")

	return int64(userId)
}

// 设置用户会话信息
func (c AppController) setLoginSession(sessionUser models.SessionUser) {
	c.Session["login"] = sessionUser.LoginName
	c.Session["ucode"] = sessionUser.Code
	c.Session["uid"] = fmt.Sprintf("%v", sessionUser.Id)
	c.Session["screen_name"] = sessionUser.DisplayName()
	c.Session["email"] = sessionUser.Email
	c.Session["from"] = sessionUser.From
	c.Session["roles"] = sessionUser.Roles
}

// 清理用户会话信息
func (c AppController) clearLoginSession() {
	delete(c.Session, "uid")
	delete(c.Session, "login")
	delete(c.Session, "ucode")
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
func (c AppController) doValidate(redirectTarget interface{}, args ...interface{}) revel.Result {
	if c.Validation.HasErrors() {
		// Store the validation errors in the flash context and redirect.
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(redirectTarget, args...)
	}
	return nil
}

func (c AppController) checkErrorAsJsonResult(err error) revel.Result {
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
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
func (c AppController) renderDTJson(page *models.PageData) revel.Result {
	sEcho := c.Params.Get("sEcho")
	return c.RenderJson(DataTableData(sEcho, page.Total, page.Total, page.Data))
}

// 构造分页查询器
func (c AppController) pageSearcher() (ps *models.PageSearcher) {
	var (
		start     int64
		limit     int64
		page      int64
		search    string
		sortColNo string
		sortField string
		sortDir   string
	)

	_, exists := c.Params.Values["iDisplayStart"]
	if exists {
		c.Params.Bind(&start, "iDisplayStart")
	} else {
		c.Params.Bind(&start, "start")
	}

	_, exists = c.Params.Values["iDisplayLength"]
	if exists {
		c.Params.Bind(&limit, "iDisplayLength")
	} else {
		c.Params.Bind(&limit, "limit")
	}

	if limit == 0 {
		limit = 10
	}

	_, exists = c.Params.Values["page"]
	if exists {
		c.Params.Bind(&page, "page")
		if page != 0 {
			start = (page - 1) * limit
		}
	} else {
		page = start/limit + 1
	}

	c.Params.Bind(&search, "sSearch")
	c.Params.Bind(&sortColNo, "iSortCol_0")
	c.Params.Bind(&sortField, "mDataProp_"+sortColNo)
	c.Params.Bind(&sortDir, "sSortDir_0")

	revel.INFO.Println("start", start, "limit", limit, "page", page)

	ps = &models.PageSearcher{
		Start: start, Limit: limit, Page: page,
		SortField: sortField, SortDir: sortDir,
		Search: search}

	ps.SetDb(c.db)

	//inject renderArgs
	c.RenderArgs["limit"] = limit
	c.RenderArgs["start"] = start
	c.RenderArgs["page"] = page
	return
}

// 构造分页查询器，通过附加的Session回调
func (c AppController) pageSearcherWithCalls(calls ...models.PageSearcherCall) *models.PageSearcher {
	ps := c.pageSearcher()
	ps.OtherCalls = calls
	return ps
}

// 设置当前页Channel
func (c AppController) setChannel(ch string) {
	c.RenderArgs["channel"] = ch
}

// 用户服务
func (c AppController) userApi() *models.UserService {
	if c.__userApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__userApi = models.NewUserService(c.db)
	}
	return c.__userApi
}

// 用户服务
func (c AppController) productApi() *models.ProductService {
	if c.__productApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__productApi = models.NewProductService(c.db)
	}
	return c.__productApi
}

func (c AppController) appApi() *models.AppService {
	if c.__appApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__appApi = models.NewAppService(c.db)
	}
	return c.__appApi
}

func (c AppController) orderApi() *models.OrderService {
	if c.__orderApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__orderApi = models.NewOrderService(c.db)
	}
	return c.__orderApi
}

func (c AppController) newsApi() *models.NewsService {
	if c.__newsApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__newsApi = models.NewNewsService(c.db)
	}
	return c.__newsApi
}

func (c AppController) appConfigApi() *models.AppConfigService {
	if c.__appConfigApi == nil {
		gotang.Assert(c.db != nil, "c.db should no be nil")
		c.__appConfigApi = models.NewAppConfigService(c.db)
	}
	return c.__appConfigApi
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
	if !isManager(c.Session) {
		return c.Redirect(App.Index)
	}
	return nil
}

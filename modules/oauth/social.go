package oauth

import (
	"encoding/base64"
	"fmt"
	//"log"
	"net/url"
	"strconv"

	"github.com/itang/yunshang/main/app/utils"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

const (
	defaultURLPrefix          = "/passport/login/"
	defaultConnectSuccessURL  = "/"
	defaultConnectFailedURL   = "/passport/login?flag=connect_failed"
	defaultLoginURL           = "/passport/login"
	defaultConnectRegisterURL = "/passport/reg/connect"
)

type SocialAuth struct {
	app                SocialAuther // 应用逻辑
	URLPrefix          string
	ConnectSuccessURL  string
	ConnectFailedURL   string
	LoginURL           string
	ConnectRegisterURL string
}

// generate session key for social-auth
// 为第三方授权生成Session Key
func (this *SocialAuth) getSessKey(social SocialType, key string) string {
	return "social_" + fmt.Sprintf("%v", social) + "_" + key
}

func (this *SocialAuth) GetSessKey(social SocialType, key string) string {
	return this.getSessKey(social, key)
}

// create oauth2 state string
// 创建oauth2 状态串
func (this *SocialAuth) createState(ctx *revel.Controller, social SocialType) string {
	values := make(url.Values, 2)

	if uid, ok := this.app.IsUserLogin(ctx); ok {
		// add uid if user current is login
		values.Add("uid", strconv.FormatInt(int64(uid), 10))
	}

	// our secret string
	values.Add("secret", utils.RandomString(20))

	// create state string
	state := base64.URLEncoding.EncodeToString([]byte(values.Encode()))

	// save to session
	name := this.getSessKey(social, "state")
	ctx.Session[name] = state

	return state
}

// verify oauth2 state string
func (this *SocialAuth) verifyState(ctx *revel.Controller, social SocialType) (string, bool) {
	code := ctx.Params.Get("code")
	state := ctx.Params.Get("state")

	if len(code) == 0 || len(state) == 0 {
		return "", false
	}

	name := this.getSessKey(social, "state")

	vu, ok := ctx.Session[name]
	if !ok || ok && state != vu {
		return "", false
	}

	return code, true
}

// Get provider according request path. ex: /login/: match /login/github
// 通过请求路径（变量）获取Provider
func (this *SocialAuth) getProvider(ctx *revel.Controller) Provider {
	path := ctx.Params.Get("provider")

	p, ok := GetProviderByPath(path)
	if !ok {
		return nil
	}

	return p
}

// After OAuthAccess check saved token for ready connect
func (this *SocialAuth) ReadyConnect(ctx *revel.Controller) (SocialType, bool) {
	var social SocialType

	sc, _ := ctx.Session["social_connect"] // 保存的类型
	s, err := strconv.Atoi(sc)
	if err != nil {
		s = 0
	}
	if s == 0 {
		return 0, false
	} else {
		social = SocialType(s)
	}

	if !social.Available() {
		return 0, false
	}

	return social, true
}

// Redirect to other social platform
// 获取重定向到第三方进行认证的链接
func (this *SocialAuth) OAuthRedirect(ctx *revel.Controller) (redirect string, failedErr error) {
	_, isLogin := this.app.IsUserLogin(ctx) // 用户是否登录， 有应用本身的逻辑判断

	defer func() {
		if len(redirect) == 0 && failedErr != nil { // unknown provider
			if isLogin {
				redirect = this.ConnectFailedURL
			} else {
				redirect = this.LoginURL
			}
		}
	}()

	var p Provider
	if p = this.getProvider(ctx); p == nil {
		failedErr = fmt.Errorf("unknown provider")
		return
	}

	social := p.GetType()
	config := p.GetConfig()
	// create redirect url
	state := this.createState(ctx, social)
	redirect = config.AuthCodeURL(state) // 认证地址
	return
}

// Callback from social platform
// 认证回调
func (this *SocialAuth) OAuthAccess(ctx *revel.Controller, session *xorm.Session) (redirect string, userSocial *UserSocial, failedErr error) {
	_, isLogin := this.app.IsUserLogin(ctx)

	defer func() {
		if len(redirect) == 0 {
			if failedErr != nil {
				if isLogin {
					redirect = this.ConnectFailedURL
				} else {
					redirect = this.LoginURL
				}
			}
		}
	}()

	// check if param has a error key
	if err := ctx.Params.Get("error"); len(err) > 0 {
		failedErr = fmt.Errorf(err)
		return
	}

	// get provider from matched url path
	var p Provider
	if p = this.getProvider(ctx); p == nil { // 请回调请求路径确定Provider
		failedErr = fmt.Errorf("unknown provider")
		return
	}

	social := p.GetType()

	var code string

	// verify state string
	// session中状态未变
	if c, ok := this.verifyState(ctx, social); !ok {
		failedErr = fmt.Errorf("state not verified")
		return
	} else {
		code = c // code 来自回调请求参数
	}

	//Cache ???
	config := p.GetConfig() // OAuth2 Config
	trans := &Transport{config, nil, DefaultTransport}

	// Send code to platform then get token
	if tok, err := trans.Exchange(code); err != nil {
		// get access token
		failedErr = err
	} else if err := tok.GetExtra("error"); err != "" {
		// token has error
		failedErr = fmt.Errorf(err)
	} else if tok.IsEmpty() {
		failedErr = fmt.Errorf("empty access token")
	} else {
		// check
		var uSocial = UserSocial{Session: session}
		if ok, err := p.CanConnect(session, tok, &uSocial); ok { // 获取tok关联的本地信息并判断是否需要Connect
			// save token to session, for connect
			// 将toke信息放到session， 方便后续进行”本地“连接
			tk := SocialTokenField{tok}
			ctx.Session[this.getSessKey(social, "token")] = tk.RawValue()
			ctx.Session["social_connect"] = strconv.Itoa(int(social))

			redirect = this.ConnectRegisterURL // 后续进入真正的“连接”处理
		} else if err == nil { // 本地已经保存了
			if !isLogin { //如果未登录
				// login user
				// 执行应用的登录逻辑
				revel.INFO.Printf("OAuthAccess: %v", uSocial)
				redirect, failedErr = this.app.LoginUser(ctx, uSocial.Uid, social)
				if len(redirect) == 0 && failedErr == nil {
					redirect = this.ConnectSuccessURL
				}
			} else {
				redirect = this.ConnectSuccessURL
			}

			// save new access token if it changed
			// 更新token到关联本地信息 如果需要的话
			uSocial.PutToken(tok)

			userSocial = &uSocial
		} else {
			failedErr = err
		}
	}

	return
}

// general use of redirect
func (this *SocialAuth) HandleRedirect(ctx *revel.Controller) revel.Result {
	redirect, err := this.OAuthRedirect(ctx)
	if err != nil {
		revel.ERROR.Printf("SocialAuth.handleRedirect, %v", err)
	}

	return ctx.Redirect(redirect)
}

// general use of redirect callback
func (this *SocialAuth) HandleAccess(ctx *revel.Controller) revel.Result {
	redirect, _, err := this.OAuthAccess(ctx, nil)
	if err != nil {
		revel.ERROR.Printf("SocialAuth.handleAccess, %v", err)
	}

	return ctx.Redirect(redirect)
}

// save user social info and login the user
// 保存
func (this *SocialAuth) ConnectAndLogin(session *xorm.Session, ctx *revel.Controller, socialType SocialType, uid int64) (string, *UserSocial, error) {
	tokKey := this.getSessKey(socialType, "token")

	defer func() {
		// delete connect tok in session
		if ctx.Session["social_connect"] != "" {
			delete(ctx.Session, "social_connect")
		}
		if ctx.Session[tokKey] != "" {
			delete(ctx.Session, tokKey)
		}
	}()

	tk := SocialTokenField{}
	value := ctx.Session[tokKey]
	if err := tk.SetRaw(value); err != nil {
		return "", nil, err
	}

	var p Provider
	if p, _ = GetProviderByType(socialType); p == nil {
		return "", nil, fmt.Errorf("unknown provider")
	}

	identify, err := p.GetIndentify(tk.Token)
	if err != nil {
		return "", nil, err
	}
	if len(identify) == 0 {
		return "", nil, fmt.Errorf("empty identify")
	}

	userSocial := UserSocial{
		Uid:      uid,
		Type:     socialType,
		RawData:  tk.String(),
		Identify: identify,
		Session:  session,
	}

	if err := userSocial.Save(); err != nil {
		return "", nil, err
	}

	// login user
	loginRedirect, err := this.app.LoginUser(ctx, uid, socialType)
	return loginRedirect, &userSocial, nil
}

// create a global SocialAuth instance
func NewSocial(urlPrefix string, socialAuther SocialAuther) *SocialAuth {
	social := new(SocialAuth)
	social.app = socialAuther

	if len(urlPrefix) == 0 {
		urlPrefix = defaultURLPrefix
	}

	if urlPrefix[len(urlPrefix)-1] != '/' {
		urlPrefix += "/"
	}

	social.URLPrefix = urlPrefix

	social.ConnectSuccessURL = defaultConnectSuccessURL
	social.ConnectFailedURL = defaultConnectFailedURL
	social.LoginURL = defaultLoginURL
	social.ConnectRegisterURL = defaultConnectRegisterURL

	return social
}

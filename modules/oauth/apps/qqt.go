package apps

import (
	//"fmt"
	//"net/url"

	//"github.com/astaxie/beego/httplib"

	"github.com/itang/yunshang/modules/oauth"
)

type QQT struct {
	BaseProvider
}

func (p *QQT) GetType() oauth.SocialType {
	return oauth.SocialQQT
}

func (p *QQT) GetName() string {
	return "QQT"
}

func (p *QQT) GetPath() string {
	return "qqt"
}

func (p *QQT) GetIndentify(tok *oauth.Token) (string, error) {
	return tok.GetExtra("openid"), nil
	/*
	uri := "https://open.t.qq.com/api/user/info?access_token=" + url.QueryEscape(tok.AccessToken)

	req := httplib.Get(uri)
	req.SetTransport(oauth.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}

	vals, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}

	if vals.Get("code") != "" {
		return "", fmt.Errorf("code: %s, msg: %s", vals.Get("code"), vals.Get("msg"))
	}

	return vals.Get("openid"), nil*/
}

var _ oauth.Provider = new(QQT)

func NewQQT(clientId, secret string) *QQT {
	p := new(QQT)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://open.t.qq.com/cgi-bin/oauth2/authorize"
	p.TokenURL = "https://open.t.qq.com/cgi-bin/oauth2/access_token"
	p.RedirectURL = oauth.DefaultAppUrl + "/passport/open/qqt/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

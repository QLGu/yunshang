package apps

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego/httplib"

	"github.com/itang/yunshang/modules/oauth"
)

type QQ struct {
	BaseProvider
}

func (p *QQ) GetType() oauth.SocialType {
	return oauth.SocialQQ
}

func (p *QQ) GetName() string {
	return "QQ"
}

func (p *QQ) GetPath() string {
	return "qq"
}

func (p *QQ) GetIndentify(tok *oauth.Token) (string, error) {
	uri := "https://graph.z.qq.com/moc2/me?access_token=" + url.QueryEscape(tok.AccessToken)

	req := httplib.Get(uri)
	//req.SetTransport(oauth.DefaultTransport)

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

	return vals.Get("openid"), nil
}

var _ oauth.Provider = new(QQ)

func NewQQ(clientId, secret string) *QQ {
	p := new(QQ)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://graph.qq.com/oauth2.0/authorize"
	p.TokenURL = "https://graph.qq.com/oauth2.0/token"
	p.RedirectURL = oauth.DefaultAppUrl + "/passport/open/qq/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

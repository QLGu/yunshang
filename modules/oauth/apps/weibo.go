package apps

import (
	"github.com/itang/yunshang/modules/oauth"
)

type Weibo struct {
	BaseProvider
}

func (p *Weibo) GetType() oauth.SocialType {
	return oauth.SocialWeibo
}

func (p *Weibo) GetName() string {
	return "Weibo"
}

func (p *Weibo) GetPath() string {
	return "weibo"
}

func (p *Weibo) GetIndentify(tok *oauth.Token) (string, error) {
	return tok.GetExtra("uid"), nil
}

var _ oauth.Provider = new(Weibo)

func NewWeibo(clientId, secret string) *Weibo {
	p := new(Weibo)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	//p.Scope = "email"
	p.AuthURL = "https://api.weibo.com/oauth2/authorize"
	p.TokenURL = "https://api.weibo.com/oauth2/access_token"
	p.RedirectURL = oauth.DefaultAppUrl + "/passport/open/weibo/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}

/**
  see https://github.com/go-av/alipay/blob/master/main.go
*/
package alipay

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
)

var alipayGatewayNew = `https://mapi.alipay.com/gateway.do?`

const (
	Service_Create_Direct_Pay_By_User = "create_direct_pay_by_user"
	DefaultPaymentType                = "1"
	BankPaymethod                     = "bankPay"
	DefaultBank                       = "" //"CMBBTB" //招行
	SuccessFeedbackCode               = "success"
	FailureFeedbackCode               = "failure"
)

type Config struct {
	Partner string // 合作身份者ID ， 以2088开头由16位纯数字组成的字符串
	Key     string //商户的私钥

	Service     string // 签约服务 create_direct_pay_by_user\必填
	PaymentType string // 支付类型， 必填，不能修改
	SellerEmail string //卖家支付宝帐户. 必填
	NotifyUrl   string //服务器异步通知页面路径， 需http://格式的完整路径，不能加?id=123这类自定义参数
	ReturnUrl   string //页面跳转同步通知页面路径， 需http://格式的完整路径，不能加?id=123这类自定义参数，不能写成http://localhost/
}

type Request struct {
	OutTradeNo string  //卖家支付宝帐户. 商户网站订单系统中唯一订单号，必填
	Subject    string  //订单名称. 必填
	TotalFee   float64 //付款金额. 必填

	Body    string //订单描述
	ShowUrl string //商品展示地址, 需以http://开头的完整路径，例如：http://www.xxx.com/myorder.html

	ExterInvokeIp string //客户端的IP地址, 非局域网的外网IP地址，如：221.0.0.1

	AntiPhishingKey string //防钓鱼时间戳

	Paymethod   string //默认支付方式. 网银支付必填
	Defaultbank string //默认网银. 网银支付必填，银行简码请参考接口技术文档
}

func NewBankPayRequest() Request {
	return Request{
		Paymethod:   BankPaymethod,
		Defaultbank: DefaultBank,
	}
}

func NewRequest() Request {
	return Request{}
}

type Response struct {
	BuyerEmail  string
	OutTradeNo  string
	TradeStatus string
	Subject     string
	TotalFee    float64
}

type kvpair struct {
	k, v string
}

type kvpairs []kvpair

func (t kvpairs) Less(i, j int) bool {
	return t[i].k < t[j].k
}

func (t kvpairs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t kvpairs) Len() int {
	return len(t)
}

func (t kvpairs) Sort() {
	sort.Sort(t)
}

func (t kvpairs) RemoveEmpty() (t2 kvpairs) {
	for _, kv := range t {
		if kv.v != "" {
			t2 = append(t2, kv)
		}
	}
	return
}

func (t kvpairs) RemoveCustom() (t2 kvpairs) {
	for _, kv := range t {
		if kv.k != "METHOD" { // METHOD , revel
			t2 = append(t2, kv)
		}
	}
	return
}

func (t kvpairs) Join() string {
	var strs []string
	for _, kv := range t {
		strs = append(strs, kv.k+"="+kv.v)
	}
	return strings.Join(strs, "&")
}

func md5Sign(str, key string) string {
	h := md5.New()
	io.WriteString(h, str)
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func verifySign(c Config, u url.Values) (err error) {
	p := kvpairs{}
	sign := ""
	for k := range u {
		v := u.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			continue
		}
		p = append(p, kvpair{k, v})
	}
	if sign == "" {
		err = fmt.Errorf("sign not found")
		return
	}
	p = p.RemoveEmpty()
	p = p.RemoveCustom() //add by itang

	p.Sort()

	if md5Sign(p.Join(), c.Key) != sign {
		err = fmt.Errorf("sign invalid")
		return
	}
	return
}

func ParseResponse(c Config, p url.Values) (r Response, err error) {
	if err = verifySign(c, p); err != nil {
		return
	}

	r.BuyerEmail = p.Get("buyer_email")
	r.TradeStatus = p.Get("trade_status")
	r.OutTradeNo = p.Get("out_trade_no")
	r.Subject = p.Get("subject")
	fmt.Sscanf(p.Get("total_fee"), "%f", &r.TotalFee)

	if r.TradeStatus != "TRADE_SUCCESS" && r.TradeStatus != "TRADE_FINISHED" {
		err = fmt.Errorf("trade not success or finnished")
		return
	}
	return
}

func NewPage(c Config, r Request, w io.Writer) {
	p := kvpairs{
		kvpair{`_input_charset`, `utf-8`},
		kvpair{`payment_type`, c.PaymentType},
		kvpair{`partner`, c.Partner},
		kvpair{`notify_url`, c.NotifyUrl},
		kvpair{`return_url`, c.ReturnUrl},
		kvpair{`service`, c.Service},
		kvpair{`seller_email`, c.SellerEmail},
		kvpair{`out_trade_no`, r.OutTradeNo},
		kvpair{`subject`, r.Subject},
		kvpair{`total_fee`, fmt.Sprintf("%.2f", r.TotalFee)},
		kvpair{`body`, r.Body},
		kvpair{`show_url`, r.ShowUrl},
		kvpair{`anti_phishing_key`, r.AntiPhishingKey},
		kvpair{`exter_invoke_ip`, r.ExterInvokeIp},
		kvpair{`paymethod`, r.Paymethod},
		kvpair{`defaultbank`, r.Defaultbank},
	}
	p = p.RemoveEmpty()
	p.Sort()

	sign := md5Sign(p.Join(), c.Key)
	p = append(p, kvpair{`sign`, sign})
	p = append(p, kvpair{`sign_type`, `MD5`})

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='post'> `, alipayGatewayNew)
	for _, kv := range p {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.k, kv.v)
	}
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
}

package entity

import (
	"time"

	"github.com/itang/yunshang/modules/alipay"
)

var (
	ACS_S  = "site"                 //站点设置
	ACS_SB = full(ACS_S, "basic")   //基本设置
	ACS_SC = full(ACS_S, "comment") //评论设置
	ACS_SO = full(ACS_S, "contact") //联系信息设置
	ACS_SM = full(ACS_S, "mail")    //邮件设置
	ACS_SP = full(ACS_S, "alipay")  //阿里支付
)

func NewAppConfig(name, key, value, description string) AppConfig {
	return AppConfig{Name: name, Key: key, Value: value, Description: description}
}

type AppConfig struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`                       //命名
	Key   string `json:"key"`                        //KEY
	Value string `xorm:"varchar(4000)" json:"value"` //VALUE

	Description string `json:"description"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (e AppConfig) IsTextArea() bool {
	switch e.Key {
	case "site.basic.statcode":
		return true
	case "site.basic.links":
		return true
	}
	return false
}

func full(s1 string, s2 string) string {
	return s1 + "." + s2
}

////////////////////////////////////////
var DefaultAppConfs = []AppConfig{
	NewAppConfig("站点名称", full(ACS_SB, "name"), "凯泰商城", "整个网站的名称"),
	NewAppConfig("默认标题", full(ACS_SB, "title"), "请修改网站名称", "整个站点默认的标题(title)，搜索引擎优化使用"),
	NewAppConfig("默认关键字", full(ACS_SB, "keywords"), "", "页面meta标签里的关键字信息(keywords)，搜索引擎优化使用"),
	NewAppConfig("默认描述信息", full(ACS_SB, "description"), "", "页面meta标签对关键字内容的描述(description)，搜索引擎优化使用"),
	NewAppConfig("ICP/IP/域名备案", full(ACS_SB, "icp"), "ICP备xxxx号", "粤ICP备00000001号"),
	NewAppConfig("版权信息", full(ACS_SB, "support"), "xxxx公司", ""),
	NewAppConfig("站点统计", full(ACS_SB, "stat"), "false", "提示：关闭后站点统计将不可用。true: 开启 false: 关闭"),
	NewAppConfig("第三方统计代码", full(ACS_SB, "statcode"), "", ""),
}

var HostAppConfs = []AppConfig{
	NewAppConfig("站点主机名", full(ACS_SB, "host"), "localhost:9000", ""),
}

var LinksAppConfs = []AppConfig{
	NewAppConfig("友情链接", full(ACS_SB, "links"), `<a href="http://www.kte99.com/" title="场效应管" target="_blank">
            场效应管</a>|
        <a href="http://kte99.cn.alibaba.com/" title="阿里巴巴" target="_blank">
            阿里巴巴</a>|
        <a href="http://www.szhxlkj.com" title="耐压测试仪" target="_blank">
            耐压测试仪</a>|
        <a href="http://www.zzsxhj.com" title="硬质合金刀片" target="_blank">
            硬质合金刀片</a>|
        <a href="http://www.dsjix.com/" title="热收缩包装机" target="_blank">
            热收缩包装机</a>|
        <a href="http://www.baoxianzuo.com" title="保险丝" target="_blank">
            保险丝</a>|
        <a href="http://www.oppower.com" title="开关电源" target="_blank">
            开关电源</a>|
        <a href="http://www.keeptops.cn" title="三防漆" target="_blank">
            三防漆</a>|
        <a href="http://www.hi1718.com/dianziyuanqijian" title="电子元器件" target="_blank">
            电子元器件</a>|
        <a href="http://luojunhong.cn.gongchang.com" title="UV胶水" target="_blank">
            UV胶水</a>|
        <a href="http://www.kte99.com/" title="场效应管" target="_blank">
            场效应管</a>|
        <a href="http://kte99.cn.alibaba.com/" title="阿里巴巴" target="_blank">
            阿里巴巴</a>|
        <a href="http://www.szhxlkj.com" title="耐压测试仪" target="_blank">
            耐压测试仪</a>|
        <a href="http://www.zzsxhj.com" title="硬质合金刀片" target="_blank">
            硬质合金刀片</a>|
        <a href="http://www.dsjix.com/" title="热收缩包装机" target="_blank">
            热收缩包装机</a>|
        <a href="http://www.baoxianzuo.com" title="保险丝" target="_blank">
            保险丝</a>|
        <a href="http://www.oppower.com" title="开关电源" target="_blank">
            开关电源</a>|
        <a href="http://www.keeptops.cn" title="三防漆" target="_blank">
            三防漆</a>|
        <a href="http://www.hi1718.com/dianziyuanqijian" title="电子元器件" target="_blank">
            电子元器件</a>|
        <a href="http://luojunhong.cn.gongchang.com" title="UV胶水" target="_blank">
            UV胶水</a>|
        <a href="http://www.kte99.com/" title="场效应管" target="_blank">
            场效应管</a>|
        <a href="http://kte99.cn.alibaba.com/" title="阿里巴巴" target="_blank">
            阿里巴巴</a>|
        <a href="http://www.szhxlkj.com" title="耐压测试仪" target="_blank">
            耐压测试仪</a>|
        <a href="http://www.zzsxhj.com" title="硬质合金刀片" target="_blank">
            硬质合金刀片</a>|
        <a href="http://www.dsjix.com/" title="热收缩包装机" target="_blank">
            热收缩包</a>|
        <a href="http://www.baoxianzuo.com" title="保险丝" target="_blank">
            保险丝</a>|`, ""),
}

var ProductAppConfs = []AppConfig{
	NewAppConfig("产品评论显示", full(ACS_SC, "show_product_comments"), "true", "true: 开启 false: 关闭"),
	NewAppConfig("新闻评论", full(ACS_SC, "show_news_comments"), "true", "true: 开启 false: 关闭"),
}

var ContactAppConfs = []AppConfig{
	NewAppConfig("公司名称", full(ACS_SO, "company_name"), "深圳市凯泰电子有限公司", ""),
	NewAppConfig("联系人", full(ACS_SO, "contact_person"), "", ""),
	NewAppConfig("联系地址", full(ACS_SO, "contact_address"), "", ""),
	NewAppConfig("联系电话", full(ACS_SO, "contact_phone"), "", ""),
	NewAppConfig("传真号码", full(ACS_SO, "fax"), "", ""),
	NewAppConfig("邮编", full(ACS_SO, "zipcode"), "", ""),
	NewAppConfig("客服邮件", full(ACS_SO, "service_email"), "", "提示：如果您有多个客服邮件，请用半角逗号（,）分隔"),
	NewAppConfig("全国销售热线", full(ACS_SO, "sales_phone"), "400-0686-198", ""),
	NewAppConfig("售后服务电话", full(ACS_SO, "after_sales_phone"), "0755-2759786", ""),
	NewAppConfig("技术支持热线", full(ACS_SO, "tech_support_phone"), "400-0686-198", ""),
	NewAppConfig("在线客服QQ", full(ACS_SO, "online_support_qq"), "售前陈R:2252410803,售前谢R:2930355581,售后王S:2252410803,售后李S:2252410803", "格式: name1:qq号1,name2:qq号2"),
	NewAppConfig("在线客服工作时间", full(ACS_SO, "online_support_time"), "8:00-17:00", ""),
}

var MoreContactAppConfs = []AppConfig{
	NewAppConfig("询价专线", full(ACS_SO, "inquiry_phone"), "0755-27597068-807", ""),
	NewAppConfig("询价QQ", full(ACS_SO, "inquiry_qq"), "2252410803", ""),
}

var MailAppConfs = []AppConfig{
	NewAppConfig("邮件服务器地址", full(ACS_SM, "host"), "smtp.gmail.com", "提示：邮件服务器地址，如：smtp.163.com，详细请查看您的邮箱设置"),
	NewAppConfig("服务器端口", full(ACS_SM, "port"), "587", "提示：邮件服务器的端口，SMTP默认为25，具体请参看各邮件服务商的设置说明"),
	NewAppConfig("加密连接(SSL)", full(ACS_SM, "ssl"), "true", "提示：此选项需要服务器环境支持SSL（如果使用Gmail或QQ邮箱，请填写true)"),
	NewAppConfig("邮件发送账号", full(ACS_SM, "username"), "yunshang2014@gmail.com", "提示：发送邮件所需的认证帐号"),
	NewAppConfig("邮件发送账号密码", full(ACS_SM, "password"), "revel2014", "提示：邮件发送账号的密码置"),
	NewAppConfig("发送者名称", full(ACS_SM, "from_name"), "凯泰商城", "提示：邮件中显示的发送者姓名"),
	NewAppConfig("发送者邮箱", full(ACS_SM, "from_address"), "yunshang2014@gmail.com", "提示：邮件中显示的发送者邮箱"),
	NewAppConfig("接收邮箱", full(ACS_SM, "recv_address"), "yunshang2014@gmail.com", "提示：用于接收系统通知邮件"),
	NewAppConfig("测试邮箱", full(ACS_SM, "test_address"), "yunshang2014@gmail.com", "提示：填写用于测试的邮件地址"),
}

var AlipayAppConfs = []AppConfig{
	NewAppConfig("合作身份者ID", full(ACS_SP, "partner"), "2088501569451024", "提示：以2088开头由16位纯数字组成的字符串"),
	NewAppConfig("商户的私钥", full(ACS_SP, "key"), "kb2hn4gmrfu3ggvexgluoqinqatids6s", "提示：商户的私钥"),
	NewAppConfig("签约服务", full(ACS_SP, "service"), alipay.Service_Create_Direct_Pay_By_User, "提示：签约服务"),
	NewAppConfig("支付类型", full(ACS_SP, "payment_type"), alipay.DefaultPaymentType, "提示：支付类型"),
	NewAppConfig("卖家支付宝帐户", full(ACS_SP, "seller_email"), "kte99168@163.com", "提示：卖家支付宝帐户email"),
	NewAppConfig("页面跳转同步通知页面路径", full(ACS_SP, "return_url"), "/payments/alipay/return", "提示：需http://格式的完整路径，不能加?id=123这类自定义参数，不能写成http://localhost/"),
	NewAppConfig("服务器异步通知页面路径", full(ACS_SP, "notify_url"), "/payments/alipay/notify", "提示：需http://格式的完整路径，不能加?id=123这类自定义参数"),
}

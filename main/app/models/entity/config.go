package entity

import (
	"time"
)

var (
	ACS_S  = "site"                 //站点设置
	ACS_SB = full(ACS_S, "basic")   //基本设置
	ACS_SC = full(ACS_S, "comment") //评论设置
	ACS_SO = full(ACS_S, "contact") //联系信息设置
)

var DefaultAppConfs = []AppConfig{
	NewAppConfig("站点名称", full(ACS_SB, "name"), "yunshang", "整个网站的名称"),
	NewAppConfig("默认标题", full(ACS_SB, "title"), "请修改网站名称", "整个站点默认的标题(title)，搜索引擎优化使用"),
	NewAppConfig("默认关键字", full(ACS_SB, "keywords"), "", "页面meta标签里的关键字信息(keywords)，搜索引擎优化使用"),
	NewAppConfig("默认描述信息", full(ACS_SB, "description"), "", "页面meta标签对关键字内容的描述(description)，搜索引擎优化使用"),
	NewAppConfig("ICP/IP/域名备案", full(ACS_SB, "icp"), "ICP备xxxx号", "粤ICP备00000001号"),
	NewAppConfig("版权信息", full(ACS_SB, "support"), "xxxx公司", ""),
	NewAppConfig("站点统计", full(ACS_SB, "stat"), "false", "提示：关闭后站点统计将不可用。true: 开启 false: 关闭"),
	NewAppConfig("第三方统计代码", full(ACS_SB, "statcode"), "", ""),
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

func NewAppConfig(name, key, value, description string) AppConfig {
	return AppConfig{Name: name, Key: key, Value: value, Description: description}
}

type AppConfig struct {
	Id int64 `json:"id"`

	Name  string `json:"name"`                       //命名
	Key   string `json:"key"`                        //KEY
	Value string `xorm:"varchar(2000)" json:"value"` //VALUE

	Description string `json:"description"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (e AppConfig) IsTextArea() bool {
	switch e.Key {
	case "site.basic.statcode":
		return true
	}
	return false
}

func full(s1 string, s2 string) string {
	return s1 + "." + s2
}

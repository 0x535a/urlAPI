package conf

var (
	Conf                 *Config
	WxToken                    = ""
	WxTokenExpireTime    int64 = -1
	MachineLogSyncLength       = 2
)

const (
	WxLoginURL        = "https://api.weixin.qq.com/sns/jscode2session"
	WxTokenURL        = "https://api.weixin.qq.com/cgi-bin/token"
	WxNotificationURL = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
	TokenLen          = 64
)

const (
	WxNotificationOrderNewTemplate     = `您有新的订单，点击查看详情。`
	WxNotificationOrderDeleteTemplate  = `订单已被删除，点击查看详情。`
	WxNotificationOrderUpdateTemplate  = `您的订单状态有变动，点击查看详情。`
	WxAppOrderDetailPageTemplate       = `/pages/order/orderDetails?orderId=%s`
	WxAppIndex                         = `/pages/index/index`
	OrderConfirmNotificationTemplateID = "-6CBB2CRgaLXasSupoffMm0cY6NPTRb2NkmjQBuzveo"
	OrderNewNotificationTemplateID     = "xHIKOyZj1nokvvo32BIzLi5HAxirikTpiNgSehwM4KY"
	OrderUpdateNotificationTemplateID  = "kRiJkmkVD91mrGCDKpFyUTQj5EnT_71R5GnWT7qZ_5c"
)

var (
	OrderStatusMap = map[string][]string{
		"Status":   {"未开始", "进行中", "待审核", "已完成"},
		"Priority": {"低", "中", "高", "紧急"},
	}
)

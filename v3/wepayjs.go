package v3

import (
	"sync"

	"github.com/valyala/fastjson"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// WepayJS 合单下单-JS支付结构体
type WepayJS struct {
	CombineAppid      string           `json:"combine_appid"`        //合单商户appid
	CombineMchid      string           `json:"combine_mchid"`        //合单商户号
	CombineOutTradeNo string           `json:"combine_out_trade_no"` //合单商户订单号
	SceneInfo         SceneInfo        `json:"scene_info"`           //场景信息
	SubOrders         SubOrders        `json:"sub_orders"`           //子单信息
	CombinePayerInfo  CombinePayerInfo `json:"combine_payer_info"`   //支付者
	TimeStart         string           `json:"time_start"`           //交易起始时间
	TimeExpire        string           `json:"time_expire"`          //交易结束时间
	NotifyURL         string           `json:"notify_url"`           //通知地址
	LimitPay          string           `json:"limit_pay"`            //指定支付方式
}

// SceneInfo 场景信息
type SceneInfo struct {
	DeviceID      string `json:"device_id"`       //商户端设备号
	PayerClientIP string `json:"payer_client_ip"` //用户终端IP
}

// SubOrders 子单信息
type SubOrders struct {
	Mchid         string     `json:"mchid"`          //子单商户号
	Attach        string     `json:"attach"`         //附加信息
	Amount        Amount     `json:"amount"`         //订单金额
	OutTradeNo    string     `json:"out_trade_no"`   //子单商户订单号
	SubMchid      string     `json:"sub_mchid"`      //二级商户号
	Detail        string     `json:"detail"`         //商品详情
	ProfitSharing bool       `json:"profit_sharing"` //是否指定分账
	Description   string     `json:"description"`    //商品描述
	SettleInfo    SettleInfo `json:"settle_info"`    //结算信息
}

// Amount 订单金额
type Amount struct {
	TotalAmount int64  `json:"total_amount"` //标价金额
	Currency    string `json:"currency"`     //标价币种
}

// SettleInfo 结算信息
type SettleInfo struct {
	ProfitSharing bool  `json:"profit_sharing"` //是否指定分账
	SubsidyAmount int64 `json:"subsidy_amount"` //补差金额
}

// CombinePayerInfo 支付者
type CombinePayerInfo struct {
	Openid string `json:"openid"` //用户标识
}

// WepayV3JSResponse 返回参数
type WepayV3JSResponse struct {
	PrepayID string `json:"prepay_id"` //预支付交易会话标识
}

var cp sync.Pool
var wrp sync.Pool

func init() {
	cp = sync.Pool{New: func() interface{} {
		return new(fasthttp.Client)
	}}
	wrp = sync.Pool{New: func() interface{} {
		return new(WepayV3JSResponse)
	}}
}
func NewWepayV3JS() *WepayJS {
	return &WepayJS{}
}
func (w *WepayJS) Do() interface{} {
	b, _ := jsoniter.Marshal(w)
	c := cp.Get().(*fasthttp.Client)
	defer cp.Put(c)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(JSAPI)
	req.SetBody(b)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := c.Do(req, resp); err != nil {
		return ""
	}
	return fastjson.GetString(resp.Body(), "prepay_id")
}

package controllers

import (
	"easy_wallpaper_api/models"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/util"
	"github.com/iGoogle-ink/gopay/pkg/xlog"
	"github.com/iGoogle-ink/gopay/wechat"
	"strconv"
	"time"
)

type WxPay struct {
	Base
}

//定义常量
//    appId：应用ID
//    mchId：商户ID
//    apiKey：API秘钥值
//    apiV3Key：apiV3Key秘钥
//    serialNo：证书序列号
const (
	AppId     = "wx7cfe2b3493f5cbc6"
	mchId     = "1511774241"
	apiKey    = "39vPDDJ4YDdvjVMRUh4fAKe93BQp9B6V"
	apiV3Key  = "8BDB05l4lVfKQSrJUWSpZgV5eXpI7xm7"
	serialNo  = "73DAC0D2BC6255926DBBF2BE0135CC6C6F75A4A7"
	pkContent = ""
)

// @Title 支付接口
// @Description 订单支付接口
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /pay_order [post]
func (this *WxPay) PayOrder() {
	// 获取支付订单id 查询支付订单详情
	sn := this.GetString("order_sn")
	data, err := models.GetOrdersPaySnInfo(sn)
	if err == nil && data.Id > 0 {
		//请求微信统一下单
		price := int64(data.TotalPrice * 100)
		sign, signErr := unifiedOrder(data.OrderSn, this.CurrentLoginUser.OpenId, this.getClientIp(), int(price))
		if signErr == nil {
			this.Data["json"] = ReturnSuccess(0, "success", sign, 1)
			this.ServeJSON()
			return
		}
		this.Data["json"] = ReturnError(40002, "微信支付请求失败")
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "无效订单,未查询到订单信息")
	this.ServeJSON()
}

type payData struct {
	Sign      string
	Package   string
	TimeStamp string
	NonceStr  string
}

// 发起微信支付请求 -- 统一下单接口
func unifiedOrder(sn, openId, ip string, price int) (payData, error) {
	//初始化微信客户端
	//    appId：应用ID
	//    mchId：商户ID
	//    apiKey：API秘钥值
	//    isProd：是否是正式环境
	client := wechat.NewClient(AppId, mchId, apiKey, true)
	//client := wechat.NewClient("wxdaa2ab9ef87b5497", "1368139502", "GFDS8j98rewnmgl45wHTt980jg543abc", false)

	//设置国家
	client.SetCountry(wechat.China)

	//number := util.GetRandomString(32)
	xlog.Debug("out_trade_no:", sn)
	//初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32))
	bm.Set("body", "小程序支付测试")
	bm.Set("out_trade_no", sn)
	bm.Set("total_fee", price)
	bm.Set("spbill_create_ip", "127.0.0.1")
	bm.Set("notify_url", "https://mp.yitiegongfang.com/v1/notify/we_chat_pay")
	bm.Set("trade_type", wechat.TradeType_Mini)
	bm.Set("device_info", "miniPro")
	bm.Set("sign_type", wechat.SignType_MD5)

	//sceneInfo := make(map[string]map[string]string)
	//miniInfo := make(map[string]string)
	//miniInfo["type"] = "Wap"
	//miniInfo["wap_url"] = "http://www.gopay.ink"
	//miniInfo["wap_name"] = "小程序测试支付"
	//sceneInfo["mini_info"] = miniInfo
	//bm.Set("scene_info", sceneInfo)

	bm.Set("openid", openId)

	// 正式
	//sign := wechat.GetParamSign("wxdaa2ab9ef87b5497", "1368139502", "GFDS8j98rewnmgl45wHTt980jg543abc", body)
	// 沙箱
	//sign, _ := wechat.GetSanBoxParamSign("wxdaa2ab9ef87b5497", "1368139502", "GFDS8j98rewnmgl45wHTt980jg543abc", body)

	// Set Sign 可以忽略不设置，内部已经自动计算sign并赋值到请求参数中了
	//bm.Set("sign", sign)
	//utils.Display("bm",bm)

	//请求支付下单，成功后得到结果
	wxRsp, err := client.UnifiedOrder(bm)
	if err != nil {
		xlog.Error(err)
		return payData{}, err
	}
	xlog.Debug("Response：", wxRsp)

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	//获取小程序支付需要的paySign
	pac := "prepay_id=" + wxRsp.PrepayId
	paySign := wechat.GetMiniPaySign(AppId, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, apiKey)
	//xlog.Debug("paySign:", paySign)
	var result payData
	result.Sign = paySign
	result.Package = pac
	result.TimeStamp = timeStamp
	result.NonceStr = wxRsp.NonceStr
	return result, nil
}

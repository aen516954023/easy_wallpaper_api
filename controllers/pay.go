package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/util"
	"github.com/iGoogle-ink/gopay/wechat/v3"
)

type Pay struct {
	Base
}

//定义常量
//    appId：应用ID
//    mchId：商户ID
//    apiKey：API秘钥值
//    isProd：是否是正式环境
const (
	appId  = "wxdaa2ab9ef87b5497"
	mchId  = "1511774241"
	apiKey = "J4sQ3YdrgAyrUznO13KKDE7e5D3j1cJz"
)

// @Title 微信支付
func (this *Pay) PayAdvanceOrder() {
	// 初始化微信客户端

	client := wechat.NewClient(appId, mchId, apiKey, false)
	// 添加微信证书 Path 路径
	//    certFilePath：apiclient_cert.pem 路径
	//    keyFilePath：apiclient_key.pem 路径
	//    pkcs12FilePath：apiclient_cert.p12 路径
	//    返回err
	client.AddCertFilePath("config/cert/apiclient_cert.pem", "config/cert/apiclient_key.pem", "config/cert/apiclient_cert.p12")

	// 初始化 BodyMap

	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32))
	bm.Set("body", "小程序测试支付")
	bm.Set("out_trade_no", "000000")
	bm.Set("total_fee", 1)
	bm.Set("spbill_create_ip", "127.0.0.1")
	bm.Set("notify_url", "http://www.gopay.ink")
	bm.Set("trade_type", gopay.TradeType_Mini)
	bm.Set("device_info", "WEB")
	bm.Set("sign_type", gopay.SignType_MD5)
	bm.Set("openid", "o0Df70H2Q0fY8JXh1aFPIRyOBgu8")

	// 嵌套json格式数据（例如：H5支付的 scene_info 参数）
	h5Info := make(map[string]string)
	h5Info["type"] = "Wap"
	h5Info["wap_url"] = "http://www.gopay.ink"
	h5Info["wap_name"] = "H5测试支付"

	sceneInfo := make(map[string]map[string]string)
	sceneInfo["h5_info"] = h5Info

	bm.Set("scene_info", sceneInfo)

	// 参数 sign ，可单独生成赋值到BodyMap中；也可不传sign参数，client内部会自动获取
	// 如需单独赋值 sign 参数，需通过下面方法，最后获取sign值并在最后赋值此参数
	sign := wechat.GetParamSign("wxdaa2ab9ef87b5497", mchId, apiKey, body)
	// sign, _ := wechat.GetSanBoxParamSign("wxdaa2ab9ef87b5497", mchId, apiKey, body)
	bm.Set("sign", sign)
}

// @Title 支付接口
// @Description 订单支付接口
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /pay_order [post]
func (this *Pay) PayOrder() {
	// 接收参数 订单id 用户id
	// 调用微信支付
	// 回调中处理，更新订单状态，更新支付状态

	orderId := this.GetString("order_id")
	payType, err := this.GetInt("pay_type")

	if err != nil {
		this.Data["json"] = ReturnError(40001, "支付单号不能为空")
		this.ServeJSON()
		this.StopRun()
	}
	// 查询订单数据
	info, errs := models.GetNotifyOrdersPay(orderId, this.CurrentLoginUser.Id)
	if errs != nil {
		this.Data["json"] = ReturnError(40001, "订单错误或订单不存在")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(info)

	switch payType {
	case 1: // 信用卡
		//postdata := make(map[string]interface{})
		//postdata["orders_code"] = info.OrderId                              //订单号
		//postdata["order_total"] = (info.TotalPrice + info.TransitPrice)     //支付总金额
		//postdata["currency_code"] = "USD"                                   //币种，例：美金USD
		//postdata["order_total_usd"] = (info.TotalPrice + info.TransitPrice) //总折算美金金额
		//postdata["notify_url"] = Config("notify_url")                       //支付结果回调地址 http://localhost:8055/notify
		//postdata["products_id"] = info.Id                                   //产品id
		//postdata["products_name"] = info.Name                               //产品名称
		//postdata["products_price"] = info.Price                             //产品价格
		//postdata["products_price_usd"] = info.Price                         //产品折算美金价格
		//data := GetOrderUrl(postdata)

		// 更新支付单号
		//boolVal, errVal := models.ModifyOrderTradeNo(info.OrderId, data["orders_id"].(string))
		//fmt.Println(boolVal)
		//if errVal == nil && boolVal {
		//	this.Data["json"] = ReturnSuccess(0, "success", data, 1)
		//	this.ServeJSON()
		//} else {
		//	logs.Error("支付请求错误:" + fmt.Sprintf("%s", errVal))
		//	this.Data["json"] = ReturnError(40003, "支付请求错误")
		//	this.ServeJSON()
		//}
		break
	default:
		this.Data["json"] = ReturnError(40004, "支付通道暂未开通")
		this.ServeJSON()
	}

}

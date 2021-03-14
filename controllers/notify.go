package controllers

import (
	"easy_wallpaper_api/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/xlog"
	"github.com/iGoogle-ink/gopay/wechat"
)

type Notify struct {
	beego.Controller
}

// @Title notify回调
// @Description 解析notify参数、验签、返回数据到微信
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /we_chat_pay [post]
func (this *Notify) ParseWeChatNotifyAndVerifyWeChatSign() {

	// 解析参数
	bodyMap, err := wechat.ParseNotifyToBodyMap(this.Ctx.Request)
	if err != nil {
		xlog.Debug("err:", err)
	}
	xlog.Debug("bodyMap:", bodyMap)

	ok, err := wechat.VerifySign(apiKey, wechat.SignType_MD5, bodyMap)
	if err != nil {
		xlog.Debug("err:", err)
	}
	xlog.Debug("微信验签是否通过:", ok)
	//2021/03/10 16:40:52.423497 /root/go/src/easy_wallpaper_api/controllers/notify.go:26: [DEBUG] >> bodyMap: map[appid:wx7cfe2b3493f5cbc6 bank_type:OTHERS cash_fee:1 device_info:miniPro fee_type:CNY is_subscribe:N mch_id:1511774241 nonce_str:PisXCvirMfTDEYTNLDHmQaknZ9OUKTfd openid:ooaG-4gBCu_4rFhJ8wq0OXNjgjfg out_trade_no:0310164016km36zb2000 result_code:SUCCESS return_code:SUCCESS sign:B518567AB7FCF9393A58A6E9DC351F0B time_end:20210310164052 total_fee:1 trade_type:JSAPI transaction_id:4200000942202103104492714898]
	//2021/03/10 16:40:52.423571 /root/go/src/easy_wallpaper_api/controllers/notify.go:32: [DEBUG] >> 微信验签是否通过: true
	//2021/03/10 16:40:52.423 [D] [server.go:2887]  |  121.51.58.169| 200 |    280.658µs|   match| POST     /v1/notify/we_chat_pay   r:/v1/notify/we_chat_pay
	//模拟支付回调
	//bodyMap :=  make(map[string]interface{})
	//bodyMap["return_code"] = "SUCCESS"
	//bodyMap["out_trade_no"] = "0310164016km36zb2000"
	//bodyMap["transaction_id"] = "4200000942202103104492714898"
	//bodyMap["time_end"] = "20210310164052"
	//ok := true

	//如果验签成功，处理业务逻辑
	if ok {
		// 更新支付状态
		status := 1
		if bodyMap["return_code"] == "SUCCESS" {
			status = 2
		} else if bodyMap["return_code"] == "FAIL" {
			status = -1
		}
		//查询支付订单信息
		payInfo, _ := models.GetOrdersPaySnInfo(bodyMap["out_trade_no"].(string))
		if payInfo.Id > 0 {
			//  更新支付订单状态 回执单号 支付时间  |  更新步骤表记录中定金支付状态
			models.UpdateOrderPayInfo(payInfo, status, bodyMap["transaction_id"].(string), bodyMap["time_end"].(string))
		} else {
			logs.Error("支付回调错误：没有查询到支付订单")
		}
	}

	rsp := new(wechat.NotifyResponse)
	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = "OK"
	//return rsp.ToXmlString()
	this.Data["xml"] = rsp
	this.ServeXML()
}

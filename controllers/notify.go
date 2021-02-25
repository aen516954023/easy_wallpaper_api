package controllers

import (
	"github.com/astaxie/beego"
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
	// 验证签名通过 处理业务逻辑
	// 1. 更新支付订单信息
	// 2. 更新订单步骤信息及状态 注意区分定金和全额的类型

	rsp := new(wechat.NotifyResponse)
	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = "OK"
	//return rsp.ToXmlString()
	this.Data["xml"] = rsp
	this.ServeXML()
}

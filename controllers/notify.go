package controllers

import (
	"github.com/astaxie/beego"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/xlog"
	"github.com/iGoogle-ink/gopay/wechat"
	"net/http"
)

type Notify struct {
	beego.Controller
}

// @Title notify回调
// @Description 解析notify参数、验签、返回数据到微信
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /we_chat_pay [post]
func (this *Notify) ParseWeChatNotifyAndVerifyWeChatSign(req *http.Request) string {
	rsp := new(wechat.NotifyResponse)
	// 解析参数
	bodyMap, err := wechat.ParseNotifyToBodyMap(req)
	if err != nil {
		xlog.Debug("err:", err)
	}
	xlog.Debug("bodyMap:", bodyMap)

	ok, err := wechat.VerifySign(apiKey, wechat.SignType_MD5, bodyMap)
	if err != nil {
		xlog.Debug("err:", err)
	}
	xlog.Debug("微信验签是否通过:", ok)

	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = "OK"
	return rsp.ToXmlString()
}

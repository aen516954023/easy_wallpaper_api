package controllers

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/util"
	"github.com/iGoogle-ink/gopay/pkg/xlog"
	"github.com/iGoogle-ink/gopay/wechat/v3"
	"io/ioutil"
	"time"
)

type Pay struct {
	beego.Controller
}

//定义常量
//    appId：应用ID
//    mchId：商户ID
//    apiKey：API秘钥值
//    apiV3Key：apiV3Key秘钥
const (
	appId     = "wxdaa2ab9ef87b5497"
	mchId     = "1511774241"
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
func (this *Pay) PayAdvanceOrder() {
	this.getPEM()
	// 初始化微信客户端
	// 	serialNo：商户证书的证书序列号
	//	pkContent：私钥 apiclient_key.pem 读取后的内容
	client, err := wechat.NewClientV3(appId, mchId, serialNo, apiV3Key, pkContent)
	if err != nil {
		xlog.Error(err)
		return
	}
	// 自动验签
	// 注意：未获取到微信平台公钥时，不要开启，请调用 client.GetPlatformCerts() 获取微信平台公钥
	//client.AutoVerifySign("微信平台公钥")

	// 打开Debug开关，输出日志
	//client.DebugSwitch = 1

	tradeNo := util.GetRandomString(32)
	xlog.Debug("tradeNo:", tradeNo)
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	// 打开Debug开关，输出日志
	client.DebugSwitch = gopay.DebugOff

	bm := make(gopay.BodyMap)
	bm.Set("description", "测试Jsapi支付商品").
		Set("out_trade_no", tradeNo).
		Set("time_expire", expire).
		Set("notify_url", "https://www.gopay.ink").
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", 1).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", "asdas")
		})

	wxRsp, err := client.V3TransactionJsapi(bm)
	if err != nil {
		xlog.Error(err)
		return
	}
	if wxRsp.Code == 0 {
		xlog.Debugf("wxRsp:%#v", wxRsp.Response)
		return
	}
	xlog.Errorf("wxRsp:%s", wxRsp.Error)
}

func (this *Pay) getPEM() {
	certFile, err := ioutil.ReadFile("conf/cert/apiclient_key.pem")
	if err != nil {
		fmt.Println(err.Error())
	}

	pemBlock, _ := pem.Decode([]byte(certFile))
	if pemBlock == nil {
		fmt.Println("decode error")
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Name %s\n", cert.Subject.CommonName)
	fmt.Printf("Not before %s\n", cert.NotBefore.String())
	fmt.Printf("Not after %s\n", cert.NotAfter.String())
}

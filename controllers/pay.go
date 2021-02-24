package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
)

/**
微信支付相关信息
商户号： 1511774241
API密钥： J4sQ3YdrgAyrUznO13KKDE7e5D3j1cJz
*/
type Pay struct {
	Base
}

// @Title 微信支付
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /pay_advance_order [post]
func (this *Pay) PayAdvanceOrder() {
	// 1 接收参数 服务订单id 支付订单id  查询支付信息
	// 2 调用微信支付
	// 3 回调中处理订单状态
	// 4 支付超时处理
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

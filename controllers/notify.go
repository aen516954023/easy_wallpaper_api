package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
	"time"
)

type Notify struct {
	beego.Controller
}

// @Title 支付回调 4111 1111 1111 1111
// @Description 支付回调接口
// @Param	orders_code		query 	string	true		"支付订单号"
// @Param	orders_status		query 	int	true		"订单状态"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /call_back [post]
func (this *Notify) CallbackNotify() {
	// orders_code 订单号  orders_status订单状态
	//$orders_status=="success" 成功状态
	//$orders_status=="failture" 失败状态

	orderCode := this.GetString("orders_code")
	orderStatus := this.GetString("orders_status")
	log.Info("订单信息:", orderCode, orderStatus)

	if orderCode == "" && orderStatus == "" {
		return
	}

	var status int
	switch orderStatus {
	case "success":
		status = 2
		break
	case "failture":
		status = 3
		break
	default:
		status = 0
	}

	// 查询订单信息
	orderInfo, err := models.GetNotifyOrdersPay(orderCode, 1)
	if err != nil {
		log.Error("查询支付单号信息错误:" + fmt.Sprintf("%s", err))
		return
	}
	// 更新 订单状态 | 支付时间
	o := orm.NewOrm()
	num, errs := o.QueryTable("order_info").
		Filter("trade_no", orderInfo.OrderSn).
		Filter("order_status", 1).
		Update(orm.Params{
			"order_status": status,
			"create_time":  time.Now().Unix(),
		})
	if errs != nil || num == 0 {
		log.Error("更新支付订单失败:" + fmt.Sprintf("%s", errs))
		return
	}
}

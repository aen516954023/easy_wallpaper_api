package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EOrdersPay))
}

type EOrdersPay struct {
	Id         int
	OrderSn    string
	OId        int
	MId        int
	TotalPrice float64
	Type       int
	PayStatus  int
	CreateAt   string
}

func GetNotifyOrdersPay(tradeNo string, mid int64) (EOrdersPay, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	err := o.QueryTable(new(EOrdersPay)).Filter("m_id", mid).Filter("trade_no", tradeNo).One(&data)
	return data, err
}

// 生成支付订单
func InsertOrderPayInfo(sn string, oId, mId, typeVal int, price float64) (int64, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	data.MId = mId
	data.OrderSn = sn
	data.OId = oId
	data.TotalPrice = price
	data.Type = typeVal
	data.PayStatus = 1
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return insertId, err
	}
	return 0, err
}

// 获取支付订单信息
func GetOrdersPayInfo(id int) (EOrdersPay, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	err := o.QueryTable("e_orders_pay").Filter("id", id).One(&data)
	return data, err
}

// 获取支付订单信息
func GetOrdersPaySnInfo(sn string) (EOrdersPay, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	err := o.QueryTable("e_orders_pay").Filter("order_sn", sn).One(&data)
	return data, err
}

// 查询订单是否存在
func GetOrderPayEmpty(oId, mId, status int) (int64, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_pay").Filter("o_id", oId).Filter("m_id", mId).Filter("pay_status", status).Count()
	return num, err
}

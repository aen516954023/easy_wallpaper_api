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
func InsertOrderPayInfo(sn string, oId, mId, typeVal int, price float64) (bool, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	data.MId = mId
	data.OrderSn = sn
	data.OId = oId
	data.TotalPrice = price
	data.Type = typeVal
	data.PayStatus = 1
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	num, err := o.Insert(&data)
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(EOrdersPay))
}

type EOrdersPay struct {
	Id      int
	OrderSn string
	OId     int
	MId     int
}

func GetNotifyOrdersPay(tradeNo string, mid int64) (EOrdersPay, error) {
	o := orm.NewOrm()
	var data EOrdersPay
	err := o.QueryTable(new(EOrdersPay)).Filter("m_id", mid).Filter("trade_no", tradeNo).One(&data)
	return data, err
}

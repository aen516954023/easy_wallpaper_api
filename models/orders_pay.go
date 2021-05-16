package models

import (
	"github.com/astaxie/beego/logs"
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
	TradeNo    string
	PayStatus  int
	PayTime    string
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
	data.Type = typeVal // 订单支付类型 1 定金  2 全额
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

//更新支付订单信息
func UpdateOrderPayInfo(d EOrdersPay, status int, tranId, payTime string) {
	o := orm.NewOrm()
	beginErr := o.Begin()
	if beginErr != nil {
		return
	}
	if d.Type == 1 { //订金
		num, err := o.QueryTable("e_orders_pay").Filter("id", d.Id).Filter("type", 1).Update(orm.Params{
			"trade_no":   tranId,
			"pay_status": status,
			"pay_time":   payTime,
		})
		if err == nil && num > 0 {
			if status == 2 {
				nums, numsErr := o.QueryTable("e_orders_step").Filter("o_id", d.OId).Update(orm.Params{
					"deposit_status": 2,
				})
				if nums > 0 && numsErr == nil {
					comErr := o.Commit()
					if comErr != nil {
						return
					}
				}
				logs.Error("更新确认基础订单支付状态失败")
			}
			comErr := o.Commit()
			if comErr != nil {
				return
			}
		}
		rollErr := o.Rollback()
		if rollErr != nil {
			return
		}

	} else if d.Type == 2 { // 全额
		num, err := o.QueryTable("e_orders_pay").Filter("id", d.Id).Filter("type", 2).Update(orm.Params{
			"trade_no":   tranId,
			"pay_status": status,
			"pay_time":   payTime,
		})
		if err == nil && num > 0 {
			if status == 2 {
				nums, numsErr := o.QueryTable("e_orders_step").Filter("o_id", d.OId).Update(orm.Params{
					"pay_status": 2,
				})
				if nums > 0 && numsErr == nil {
					comErr := o.Commit()
					if comErr != nil {
						return
					}
				}
				logs.Error("更新确认实际报价订单支付状态失败")
			}
			comErr := o.Commit()
			if comErr != nil {
				return
			}
		}
		rollErr := o.Rollback()
		if rollErr != nil {
			return
		}
	}

}

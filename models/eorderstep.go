package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EOrderStep))
}

type EOrderStep struct {
	Id               int
	OId              int
	MId              int
	WId              int
	ServiceType      int
	ConstructionType int
	Unit             int
	Price            float64
	DepositPrice     float64
	Info             string
	Status           int
	CreateAt         string
}

//添加参与记录,并更新选中的师傅状态
func InsertOrdersStep(oId, mId, wId int) bool {
	o := orm.NewOrm()
	var data EOrderStep
	beginErr := o.Begin()
	if beginErr != nil {
		logs.Error("start the transaction failed")
		return false
	}
	data.MId = mId
	data.OId = oId
	data.WId = wId
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	InsertId, err := o.Insert(&data)
	if err == nil && InsertId > 0 {
		// 更新师傅参与表中其它未选中的师傅的状态
		num, upErr := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
			"status": 1,
		})
		if upErr == nil && num > 0 {
			comErr := o.Commit()
			if comErr != nil {
				logs.Error("Commit the transaction failed")
				return false
			}
			return true
		}
		backErr := o.Rollback()
		if backErr != nil {
			logs.Error("Commit the transaction failed")
			return false
		}
		return false
	}
	backErr := o.Rollback()
	if backErr != nil {
		logs.Error("Commit the transaction failed")
		return false
	}
	return false
}

//基础报价
func InsertOrdersStepTwo(oId, mId, wId, sType, cType, unit int, price, dPrice float64, info string) bool {
	o := orm.NewOrm()
	var data EOrderStep
	beginErr := o.Begin()
	if beginErr != nil {
		logs.Error("start the transaction failed")
		return false
	}
	data.MId = mId
	data.OId = oId
	data.WId = wId
	data.ServiceType = sType
	data.ConstructionType = cType
	data.Unit = unit
	data.Price = price
	data.DepositPrice = dPrice
	data.Info = info
	data.Status = 1
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	InsertId, err := o.Insert(&data)
	if err == nil && InsertId > 0 {
		// 更新订单相关字段信息
		num, upErr := o.QueryTable("e_orders").Filter("o_id", oId).Update(orm.Params{
			"service_id":        sType,
			"construction_type": cType,
			"status":            2,
		})
		if upErr == nil && num > 0 {
			comErr := o.Commit()
			if comErr != nil {
				logs.Error("Commit the transaction failed")
				return false
			}
			return true
		}
		backErr := o.Rollback()
		if backErr != nil {
			logs.Error("Commit the transaction failed")
			return false
		}
		return false
	}
	backErr := o.Rollback()
	if backErr != nil {
		logs.Error("Commit the transaction failed")
		return false
	}
	return false
}

// 查询订单及基础报价信息
func GetOrderOfStepInfo(orderId, status int) (EOrderStep, error) {
	o := orm.NewOrm()
	var data EOrderStep
	err := o.QueryTable("e_orders_step").Filter("o_id", orderId).Filter("status", status).One(&data)
	//err := o.Raw("SELECT * FROM e_orders_step  WHERE o_id=? AND status=?",orderId,status).QueryRow(&data)
	return data, err
}

//实际报价
func InsertOrdersStepThree(oId, mId, wId, sType, cType, unit int, price, dPrice float64, info string) bool {
	o := orm.NewOrm()
	var data EOrderStep
	data.MId = mId
	data.OId = oId
	data.WId = wId
	data.ServiceType = sType
	data.ConstructionType = cType
	data.Unit = unit
	data.Price = price
	data.DepositPrice = dPrice
	data.Info = info
	data.Status = 2
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	InsertId, err := o.Insert(&data)
	if err == nil && InsertId > 0 {
		return true
	}
	return false
}

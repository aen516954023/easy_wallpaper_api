package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EOrdersStep))
}

type EOrdersStep struct {
	Id               int
	OId              int
	MId              int
	WId              int
	ServiceType      int
	ConstructionType int
	Unit             int
	Price            float64
	DepositPrice     float64
	TotalPrice       float64
	Info             string
	Status           int
	PayStatus        int
	PayId            int
	CreateAt         string
}

//添加参与记录,并更新选中的师傅状态
func InsertOrdersStep(oId, mId, wId int) bool {
	o := orm.NewOrm()
	var data EOrdersStep
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
		fmt.Println(err, upErr, num)

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
	var data EOrdersStep
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
		num, upErr := o.QueryTable("e_orders").Filter("id", oId).Update(orm.Params{
			"service_id":        sType,
			"construction_type": cType,
			"status":            2,
		})
		fmt.Println("upErr", upErr)
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
func GetOrderOfStepInfo(orderId, status int) (EOrdersStep, error) {
	o := orm.NewOrm()
	var data EOrdersStep
	err := o.QueryTable("e_orders_step").Filter("o_id", orderId).Filter("status", status).One(&data)
	//err := o.Raw("SELECT * FROM e_orders_step  WHERE o_id=? AND status=?",orderId,status).QueryRow(&data)
	return data, err
}

//实际报价
func InsertOrdersStepThree(oId, mId, wId, sType, cType, unit int, price, dPrice float64, info string) bool {
	o := orm.NewOrm()
	var data EOrdersStep
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

// 更新订单状态
func ModifyStepStatus(oId, mId, status int) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("m_id", mId).Update(orm.Params{
		"status": status,
	})
	if err == nil && num > 0 {
		return true
	}
	return false
}

func GetOrderOfStepAll(orderId, mId int) (int64, []EOrdersStep, error) {
	o := orm.NewOrm()
	var data []EOrdersStep
	num, err := o.QueryTable("e_orders_step").Filter("o_id", orderId).Filter("m_id", mId).OrderBy("-create_at").All(&data)
	return num, data, err
}

// 更新订单步骤支付状态
func UpdateOrderStepPayStatus(orderId int, mId int64, status, payStatue int, insertId int64) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", orderId).
		Filter("m_id", mId).Filter("status", status).Update(orm.Params{
		"pay_status": payStatue,
		"pay_id":     insertId,
	})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

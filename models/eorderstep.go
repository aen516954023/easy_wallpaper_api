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
	Area             float64
	Price            float64
	DiscountedPrice  float64
	Info             string
	TotalPrice       float64
	Unit             int
	DepositPrice     float64
	HomeTime         int
	Step1            int
	Step1Time        int
	Step2            int
	Step2Time        int
	Step3            int
	Step3Time        int
	Step4            int
	Step4Time        int
	Step5            int
	Step5Time        int
	Step6            int
	Step6Time        int
	Step7            int
	Step7Time        int
	DepositStatus    int
	DepositId        int
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
		// 更新师傅参与表中选中的师傅状态
		num, upErr := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
			"status": 1,
		})
		fmt.Println(err, upErr, num)

		if upErr == nil && num > 0 {
			// 更新订单表中状态== 2  进入订单流程
			ordersNum, eOrdersErr := o.QueryTable("e_orders").Filter("id", oId).Filter("m_id", mId).Update(orm.Params{
				"status": 2,
			})
			if eOrdersErr == nil && ordersNum > 0 {
				comErr := o.Commit()
				if comErr != nil {
					logs.Error("Commit the transaction failed")
					return false
				}
				return true
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
	backErr := o.Rollback()
	if backErr != nil {
		logs.Error("Commit the transaction failed")
		return false
	}
	return false
}

//基础报价
func ModifyOrdersStepTwo(oId, wId, sType, cType, unit int, price, dPrice float64, info string, hTime int64) bool {
	o := orm.NewOrm()
	beginErr := o.Begin()
	if beginErr != nil {
		logs.Error("start the transaction failed")
		return false
	}
	nums, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
		"service_type":      sType,
		"construction_type": cType,
		"price":             price,
		"deposit_price":     dPrice,
		"unit":              unit,
		"info":              info,
		"home_time":         hTime,
		"step1":             1,
		"step1_time":        time.Now().Unix(),
	})
	//data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	if err == nil && nums > 0 {
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
func GetOrderOfStepInfo(orderId int) (EOrdersStep, error) {
	o := orm.NewOrm()
	var data EOrdersStep
	err := o.QueryTable("e_orders_step").Filter("o_id", orderId).One(&data)
	return data, err
}

//实际报价
func ModifyOrdersStepActualQuotation(oId, wId, sType, cType, unit int, area, price, dPrice, tPrice float64, info string) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
		"service_type":      sType,
		"construction_type": cType,
		"area":              area,
		"price":             price,
		"discounted_price":  dPrice,
		"total_price":       tPrice,
		"unit":              unit,
		"info":              info,
		"step4":             1,
		"step4_time":        time.Now().Unix(),
	})
	if err == nil && num > 0 {
		return true
	}
	return false
}

// 发起验收
func ModifyStepStatus(oId, wId int) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
		"step6":      1,
		"step6_time": time.Now().Unix(),
	})
	if err == nil && num > 0 {
		return true
	}
	return false
}

// 确认验收
func ModifyConfirmAcceptance(oId, mId int) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("m_id", mId).Update(orm.Params{
		"step7":      1,
		"step7_time": time.Now().Unix(),
	})
	if err == nil && num > 0 {
		return true
	}
	return false
}

func GetOrderOfStepOne(orderId, mId int) (EOrdersStep, error) {
	o := orm.NewOrm()
	var data EOrdersStep
	err := o.QueryTable("e_orders_step").Filter("o_id", orderId).Filter("m_id", mId).OrderBy("-create_at").One(&data)
	return data, err
}

// 更新用户基础报价确认步骤状态
func UpdateOrderStep2(orderId int, mId, insertId int64) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", orderId).
		Filter("m_id", mId).Update(orm.Params{
		"step2":          1,
		"step2_time":     time.Now().Unix(),
		"deposit_id":     insertId,
		"deposit_status": 1,
	})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

// 更新师傅已到现场状态
func ModifyStep3Status(oId, wId int) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
		"step3":      1,
		"step3_time": time.Now().Unix(),
	})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

// 更新用户基础报价确认步骤状态
func UpdateOrderStep5(orderId int, mId, insertId int64) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders_step").Filter("o_id", orderId).
		Filter("m_id", mId).Update(orm.Params{
		"step5":      1,
		"step5_time": time.Now().Unix(),
		"pay_id":     insertId,
		"pay_status": 1,
	})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

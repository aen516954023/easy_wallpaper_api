package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type EMemberOrMasterWorker struct {
	Id       int
	OId      int
	MId      int
	WId      int
	Status   int
	CreateAt string
}

func init() {
	orm.RegisterModel(new(EMemberOrMasterWorker))
}

//查询订单参与次数
func GetTaskCount(oId int) (int64, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Count()
	return num, err
}

//查询师傅接单次数
func GetOrderTasking(oId, wId int) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Filter("w_id", wId).Count()
	if err == nil && num > 0 {
		return true
	}
	return false
}

//查询师傅是否参与订单 uid查询
func GetOrderTaskingUid(oId, uid int) bool {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Filter("m_id", uid).Count()
	if err == nil && num > 0 {
		return true
	}
	return false
}

// 师傅接单表写入数据
func InsertOrderTaking(oid, mId, wid int) (bool, error) {
	o := orm.NewOrm()
	var data EMemberOrMasterWorker
	data.OId = oid
	data.MId = mId
	data.WId = wid
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

type OrderWorkerList struct {
	Id           int
	Name         string
	Image        string
	ServiceCount int
	IsRealName   int
	Exp          int
	Warranty     int
}

//参与报价的师傅列表数据
func GetOrderWorkerList(orderId int) (int64, []OrderWorkerList, error) {
	o := orm.NewOrm()
	var data []OrderWorkerList
	num, err := o.Raw("SELECT w.w_id as id,m.username as name,m.image,m.is_real_name,m.exp,m.warranty FROM e_member_or_master_worker w LEFT JOIN e_master_worker m ON m.id=w.w_id WHERE w.o_id=? and w.status >= 0", orderId).QueryRows(&data)
	return num, data, err
}

//师傅参与的订单
type OrderMasterList struct {
	Id               int
	OrderSn          string
	ServiceId        int
	ConstructionType int
	ConstructionTime int
	Area             float64
	CreateAt         string
}

func GetOrderMasterAll(mId int64, status int) (int64, []OrderMasterList, error) {
	o := orm.NewOrm()
	var data []OrderMasterList
	if status == 0 {
		num, err := o.Raw("SELECT o.id,o.order_sn,o.service_id,o.construction_type,o.construction_time,o.area,o.create_at FROM e_member_or_master_worker w LEFT JOIN e_orders o ON o.id=w.o_id WHERE w.m_id=? and w.status=1 and o.status > ?", mId, status).QueryRows(&data)
		return num, data, err
	} else if status < 0 {
		num, err := o.Raw("SELECT o.id,o.order_sn,o.service_id,o.construction_type,o.construction_time,o.area,o.create_at FROM e_member_or_master_worker w LEFT JOIN e_orders o ON o.id=w.o_id WHERE w.m_id=? and w.status=0", mId).QueryRows(&data)
		return num, data, err
	} else {
		num, err := o.Raw("SELECT o.id,o.order_sn,o.service_id,o.construction_type,o.construction_time,o.area,o.create_at FROM e_member_or_master_worker w LEFT JOIN e_orders o ON o.id=w.o_id WHERE w.m_id=? and w.status=1 and o.status = ?", mId, status).QueryRows(&data)
		return num, data, err
	}

}

// 淘汰师傅 更新状态
func CancelMasterOrder(oId, wId int) (bool, error) {
	o := orm.NewOrm()
	// 更新师傅参与表状态
	num, err := o.QueryTable("e_member_or_master_worker").Filter("o_id", oId).Filter("w_id", wId).Update(orm.Params{
		"status": -1,
	})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

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

// 师傅接单表写入数据
func InsertOrderTaking(oid, wid int) (bool, error) {
	o := orm.NewOrm()
	var data EMemberOrMasterWorker
	data.OId = oid
	data.MId = 0
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
}

//参与报价的师傅列表数据
func GetOrderWorkerList(orderId int) (int64, []OrderWorkerList, error) {
	o := orm.NewOrm()
	var data []OrderWorkerList
	num, err := o.Raw("SELECT * FROM e_member_of_master_worker w LEFT JOIN e_master_worker m ON m.id=w.w_id").QueryRows(&data)
	return num, data, err
}

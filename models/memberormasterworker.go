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

// 师傅接单表写入数据
func InsertOrderTaking(oid, mid, wid int) bool {
	o := orm.NewOrm()
	var data EMemberOrMasterWorker
	data.OId = oid
	data.MId = mid
	data.WId = wid
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true
	}
	return false
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

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

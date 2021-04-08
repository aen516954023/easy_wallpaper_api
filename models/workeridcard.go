package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EMasterWorkerIdcard))
}

type EMasterWorkerIdcard struct {
	Id       int
	WId      int
	IdCard   string
	RealName string
	ImagePos string
	ImageNeg string
	Status   int
	CreateAt string
}

//实名认证提交
func AddIdCard(wId int, idCard, realname, pos, neg string) (bool, error) {
	o := orm.NewOrm()
	var data EMasterWorkerIdcard
	data.WId = wId
	data.IdCard = idCard
	data.RealName = realname
	data.ImagePos = pos
	data.ImageNeg = neg
	data.Status = 0
	data.CreateAt = time.Now().Format("2016-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

// 获取实名认证信息
func GetIdCard(wId int) (EMasterWorkerIdcard, error) {
	o := orm.NewOrm()
	var data EMasterWorkerIdcard
	err := o.QueryTable("e_master_worker_idcard").Filter("w_id", wId).One(&data)
	return data, err
}

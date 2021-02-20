package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(EMasterWorker))
}

type EMasterWorker struct {
	Id          int
	MId         int
	Username    string
	Mobile      string
	Gender      int
	Exp         int
	Warranty    int
	Address     string
	ServiceCity string
	Describe    string
	Image       string
	IsRealName  int
	Status      int
	CreateAt    string
}

func GetMasterWorkerInfo(uid int64) (EMasterWorker, error) {
	o := orm.NewOrm()
	var data EMasterWorker
	err := o.QueryTable("e_master_worker").Filter("m_id", uid).One(&data)
	return data, err
}

//获取师傅列表
func GetRecommendList(limit int) (int64, []EMasterWorker, error) {
	o := orm.NewOrm()
	var data []EMasterWorker
	num, err := o.QueryTable("e_master_worker").Filter("status", 1).Limit(limit).All(&data)
	return num, data, err
}

// 师傅入驻申请
func ApplyMasterWorker(mid int64, gender, exp int, userName, mobile, city, address, desc, avatar, addTime string) (bool, error) {
	o := orm.NewOrm()
	var data EMasterWorker
	data.MId = int(mid)
	data.Username = userName
	data.Gender = gender
	data.Mobile = mobile
	data.Exp = exp
	data.ServiceCity = city
	data.Address = address
	data.Describe = desc
	data.Image = avatar
	data.CreateAt = addTime

	num, err := o.Insert(&data)
	if err == nil && num > 0 {
		return true, err
	}
	return false, err
}

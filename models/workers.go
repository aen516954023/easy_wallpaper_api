package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(EMasterWorker))
}

type EMasterWorker struct {
	Id         int
	Mid        int
	Image      string
	IsRealName int
}

func GetRecommendList(limit int) (int64, []EMasterWorker, error) {
	o := orm.NewOrm()
	var data []EMasterWorker
	num, err := o.QueryTable(new(EMasterWorker)).Filter("status", 1).Limit(limit).All(&data)
	return num, data, err
}

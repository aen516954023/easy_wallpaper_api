package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EMasterWorkerExp))
}

type EMasterWorkerExp struct {
	Id       int
	Wid      int
	Exp      int
	Phone    string
	Wechat   string
	Image    string
	Status   int
	CreateAt string
}

//获取师傅经验信息
func GetExp(wId int) (*EMasterWorkerExp, error) {
	o := orm.NewOrm()
	var data *EMasterWorkerExp
	err := o.QueryTable("e_master_worker_exp").Filter("w_id", wId).One(&data)
	return data, err
}

func AddExp(wId, exp int, wechat, phone, image string) (bool, error) {
	o := orm.NewOrm()
	var data EMasterWorkerExp
	data.Wid = wId
	data.Exp = exp
	data.Phone = phone
	data.Image = image
	data.Wechat = wechat
	data.Status = 0
	data.CreateAt = time.Now().Format("2016-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EMasterWorkerCase))
}

type EMasterWorkerCase struct {
	Id       int
	WId      int
	Content  string
	Images   string
	CreateAt string
}

//添加案例
func AddMasterCase(wid int, desc, images string) (bool, error) {
	o := orm.NewOrm()
	var data EMasterWorkerCase
	data.WId = wid
	data.Content = desc
	data.Images = images
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

//案例列表
func GetMasterCaseList(wid, page, pageNum int) (int64, []EMasterWorkerCase, error) {
	o := orm.NewOrm()
	var data []EMasterWorkerCase
	num, err := o.Raw("select id,w_id,content, images, create_at from e_master_worker_case where w_id=? and id > (?-1)*?  limit ?", wid, page, pageNum, pageNum).QueryRows(&data)
	return num, data, err
}

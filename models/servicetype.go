package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(EServiceType))
}

type EServiceType struct {
	Id       int
	TypeName string
	Status   int
	CreateAt string
}

func GetAllServiceType() (int64, []EServiceType, error) {
	o := orm.NewOrm()
	var data []EServiceType
	num, err := o.QueryTable(new(EServiceType)).All(&data)
	return num, data, err
}

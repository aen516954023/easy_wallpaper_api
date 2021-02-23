package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(EOrderStep))
}

type EOrderStep struct {
	Id               int
	OId              int
	MId              int
	WId              int
	ServiceType      int
	ConstructionType int
	Status           int
	CreateAt         string
}

//添加参与记录
func InsertMasterOrder(oId, mId, wId int, cTime string) bool {
	o := orm.NewOrm()
	var data EOrderStep
	data.MId = mId
	data.OId = oId
	data.WId = wId
	data.ServiceType = 1
	data.ConstructionType = 1
	data.CreateAt = cTime
	InsertId, err := o.Insert(&data)
	if err == nil && InsertId > 0 {
		return true
	}
	return false
}

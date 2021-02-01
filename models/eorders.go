package models

import "github.com/astaxie/beego/orm"

type EOrders struct {
	Id                   int
	MId                  int
	WorkerId             int
	Address              string
	ConstructionTime     int
	ServiceId            int
	IsMateriel           int
	Area                 float64
	IsTearOfOldWallpaper int
	BasementMembrane     int
	MoreDescription      string
	OrderType            int
	Status               int
	CreateAt             string
}

func init() {
	orm.RegisterModel(new(EOrders))
}

//查看所有订单列表
func GetOrdersAll() (int64, []EOrders, error) {
	o := orm.NewOrm()
	var data []EOrders
	num, err := o.QueryTable("e_orders").All(&data)
	return num, data, err
}

package models

import "github.com/astaxie/beego/orm"

type EOrders struct {
	Id                   int
	MId                  int
	WorkerId             int
	Address              string
	ConstructionTime     int
	ConstructionType     int
	ServiceId            int
	IsMateriel           int
	Area                 float64
	IsTearOfOldWallpaper int
	BasementMembrane     int
	MoreDescription      string
	OrderType            int
	Status               int
	Images               string
	CreateAt             string
}

func init() {
	orm.RegisterModel(new(EOrders))
}

//查看所有订单列表
func GetOrdersAll(status int) (int64, []EOrders, error) {
	o := orm.NewOrm()
	var data []EOrders
	if status == 0 {
		num, err := o.QueryTable("e_orders").All(&data)
		return num, data, err
	} else {
		num, err := o.QueryTable("e_orders").Filter("status", status).All(&data)
		return num, data, err
	}
}

// 查询订单详情
func GetOrderInfo(orderId int) (EOrders, error) {
	o := orm.NewOrm()
	var data EOrders
	err := o.QueryTable("e_orders").Filter("id", orderId).One(&data)
	return data, err
}

//取消订单
func OrderCancel(orderId, s int) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("e_orders").Filter("id", orderId).Update(orm.Params{
		"status": s,
	})
	if num > 0 && err == nil {
		return true, err
	}
	return false, err
}

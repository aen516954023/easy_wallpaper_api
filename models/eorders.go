package models

import "github.com/astaxie/beego/orm"

type EOrders struct {
	Id                   int
	OrderSn              string
	MId                  int
	WorkerId             int
	Address              int
	City                 string
	ConstructionTime     int
	ConstructionTimeStr  string
	ServiceId            int
	ConstructionType     int
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
func GetOrdersAll(mId int64, status, flag int) (int64, []EOrders, error) {
	o := orm.NewOrm()
	var data []EOrders
	// flag 1 用户端订单列表  2 师傅端 订单大厅
	if flag == 2 {
		// Exclude 过滤当前师傅用户发布的订单 防止师傅自己接自己的订单
		if status == 0 {
			num, err := o.QueryTable("e_orders").Exclude("m_id", mId).Filter("status__gt", 0).OrderBy("-create_at").All(&data)
			return num, data, err
		} else {
			// status = 1 进行中的订单
			// 获取当前用户师傅信息服务的城市
			workerInfo, _ := GetMasterWorkerInfo(mId)

			num, err := o.QueryTable("e_orders").Exclude("m_id", mId).Filter("status", status).Filter("city", workerInfo.ServiceCity).OrderBy("-create_at").All(&data)
			return num, data, err
		}
	} else {
		if status == 0 {
			num, err := o.QueryTable("e_orders").Filter("m_id", mId).Filter("status__gt", 0).OrderBy("-create_at").All(&data)
			return num, data, err
		} else {
			num, err := o.QueryTable("e_orders").Filter("m_id", mId).Filter("status", status).OrderBy("-create_at").All(&data)
			return num, data, err
		}
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

type UserOpenId struct {
	OpenId   string
	OrderSn  string
	CreateAt string
}

//获取订单用户openid
func GetOrderUserOpenid(orderId int) (UserOpenId, error) {
	o := orm.NewOrm()
	var data UserOpenId
	err := o.Raw("SELECT m.open_id,o.order_sn,o.create_at FROM e_orders o LEFT JOIN e_members m ON m.id=o.m_id WHERE o.id=?", orderId).QueryRow(&data)
	return data, err
}

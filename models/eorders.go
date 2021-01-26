package models

import "github.com/astaxie/beego/orm"

type EOrders struct {
	Id                   int
	OrderSn              string
	Mid                  int
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

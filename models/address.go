package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(EAddress))
}

// 地址表
type EAddress struct {
	Id        int
	MId       int     // 用户id
	Username  string  // 联系人
	Phone     string  // 联系方式
	Name      string  // 地址名称
	Province  string  // 省
	City      string  // 市
	District  string  // 区
	Address   string  // 地址详情
	Latitude  float64 // 经度
	Longitude float64 // 纬度
	Default   int     // 是否是默认地址
	CreateAt  string  //添加时间
}

func GetAddressList(mId int64) (int64, []EAddress, error) {
	o := orm.NewOrm()
	var data []EAddress
	num, err := o.QueryTable("e_address").Filter("m_id", mId).All(&data)
	return num, data, err
}

func GetAddressId(id int) (EAddress, error) {
	o := orm.NewOrm()
	var data EAddress
	err := o.QueryTable("e_address").Filter("id", id).One(&data)
	return data, err
}

func InsertAddress(insert interface{}) (bool, error) {
	//o := orm.NewOrm()
	//var data EAddress
	//data.MId = insert.Id
	return true, nil
}

package models

import (
	"github.com/astaxie/beego/orm"
	"time"
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

//添加操作
func InsertAddress(mId, isDefault int, username, phone, addressName, address, pro, city, dist string, lat, log float64) (int64, error) {
	o := orm.NewOrm()
	o.Begin()
	// 判断添加地址是否设置默认， 如果设置 取消其它默认地址
	if isDefault > 0 {
		o.QueryTable("e_address").Filter("m_id", mId).Update(orm.Params{
			"default": 0,
		})
	}

	var data EAddress
	data.MId = mId
	data.Username = username
	data.Phone = phone
	data.Name = addressName
	data.Province = pro
	data.City = city
	data.District = dist
	data.Address = address
	data.Latitude = lat
	data.Longitude = log
	data.Default = isDefault
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		o.Commit()
		return insertId, err
	}
	o.Rollback()
	return 0, err
}

// 修改操作
func ModifyAddress(id, mId, isDefault int, username, phone, addressName, address, pro, city, dist string, lat, log float64) (bool, error) {
	o := orm.NewOrm()
	o.Begin()
	// 判断添加地址是否设置默认， 如果设置 取消其它默认地址
	if isDefault > 0 {
		o.QueryTable("e_address").Filter("m_id", mId).Update(orm.Params{
			"default": 0,
		})
	}

	num, err := o.QueryTable("e_address").Filter("id", id).Update(orm.Params{
		"username":  username,
		"phone":     phone,
		"name":      addressName,
		"province":  pro,
		"city":      city,
		"district":  dist,
		"address":   address,
		"Latitude":  lat,
		"Longitude": log,
		"default":   isDefault,
	})
	if err == nil && num > 0 {
		o.Commit()
		return true, err
	}
	o.Rollback()
	return false, err
}

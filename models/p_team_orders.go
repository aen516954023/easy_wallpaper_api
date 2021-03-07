package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(PTeamOrders))
}

type PTeamOrders struct {
	Id     int
	Name   string
	Gander int
	Area   float64
	Day    int
	City   string
	Phone  string
}

func SaveTeamOrders(name, phone, city string, gender, day int, area float64) (bool, error) {
	o := orm.NewOrm()
	var data = PTeamOrders{}
	data.Name = name
	data.Gander = gender
	data.Area = area
	data.City = city
	data.Phone = phone
	data.Day = day
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

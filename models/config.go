package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Config))
}

type Config struct {
	Id     int
	Fields string
	Value  string
}

func GetConfig(field string) (Config, error) {
	var data Config
	o := orm.NewOrm()
	err := o.QueryTable("config").Filter("fields", field).One(&data)
	return data, err
}

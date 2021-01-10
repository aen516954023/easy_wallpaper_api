package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Article))
}

// 地址表
type Article struct {
	Id       int
	Title    string
	Desc     string
	Content  string
	Type     int
	CreateAt string
}

func GetMsgBefore(nums int) (int64, []Article, error) {
	var data []Article
	o := orm.NewOrm()
	num, err := o.QueryTable(new(Article)).Limit(nums).All(&data)
	return num, data, err
}

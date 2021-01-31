package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(EMembers))
}

type EMembers struct {
	Id     int64
	OpenId string
}

//通过openid查询用户信息
func GetMemberInfo(openId string) (EMembers, error) {
	o := orm.NewOrm()
	var data EMembers
	err := o.QueryTable("e_members").Filter("open_id", openId).One(&data)
	return data, err
}

//创建新用户
func AddMember(openId string) (int64, error) {
	o := orm.NewOrm()
	var data EMembers
	data.OpenId = openId
	return o.Insert(&data)
}

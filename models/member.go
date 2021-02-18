package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(EMembers))
}

type EMembers struct {
	Id       int64
	Nickname string
	OpenId   string
	Phone    string
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

// 更新用户手机号码
func UpdatePhone(id int64, phone string) bool {
	o := orm.NewOrm()
	id, err := o.QueryTable("e_members").Filter("id", id).Update(orm.Params{
		"phone": phone,
	})
	fmt.Println(err)
	if err == nil && id > 0 {
		return true
	}
	return false
}

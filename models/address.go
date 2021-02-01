package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Address))
}

// 地址表
type Address struct {
	Id        int
	Receiver  string       `orm:"size(20)"`       // 收件人
	Addr      string       `orm:"size(50)"`       // 收件地址
	ZipCode   string       `orm:"size(20)"`       // 邮编
	Phone     string       `orm:"size(20)"`       // 联系方式
	IsDefault bool         `orm:"default(false)"` // 是否是默认地址
	User      *EMembers    `orm:"rel(fk)"`        // 用户ID
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

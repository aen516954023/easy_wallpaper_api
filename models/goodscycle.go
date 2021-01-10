package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(GoodsCycle))
}

// 商品类型表
type GoodsCycle struct {
	Id       int
	Day      int         // 周期天数
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"` // 商品SKU
}

func GetAllzq() (int64, []GoodsCycle, error) {
	o := orm.NewOrm()
	var info []GoodsCycle
	num, err := o.Raw("SELECT * FROM goods_cycle ").QueryRows(&info)
	return num, info, err
}

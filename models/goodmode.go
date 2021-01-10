package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(GoodsMode))
}

// 商品类型表
type GoodsMode struct {
	Id       int
	Name     string      `orm:"size(20)"`      // 类型名称
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"` // 商品SKU
}

func GetAlltc() (int64, []GoodsMode, error) {
	o := orm.NewOrm()
	var info []GoodsMode
	num, err := o.Raw("SELECT * FROM goods_mode ").QueryRows(&info)
	return num, info, err
}

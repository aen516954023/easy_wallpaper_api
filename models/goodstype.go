package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(GoodsType))
}

// 商品类型表
type GoodsType struct {
	Id       int
	Name     string      `orm:"size(20)"` // 类型名称
	Logo     string      // 类型Logo
	Price    float64     // 币种价格
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"` // 商品SKU
}

// 获取所有币种信息
func GetAllCurrency() (int64, []GoodsType, error) {
	o := orm.NewOrm()
	var info []GoodsType
	num, err := o.Raw("SELECT * FROM goods_type ").QueryRows(&info)
	return num, info, err
}

func Get_first_type_id() (GoodsType, error) {
	o := orm.NewOrm()
	var info GoodsType
	err := o.Raw("SELECT * from goods_type limit 1").QueryRow(&info)
	return info, err
}

package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Goods))
}

// 商品SPU表
type Goods struct {
	Id       int
	Name     string      `orm:"size(2000)"` // 矿机名称
	Detail   string      `orm:"size(200)"`  // 详细描述
	Hashrate string      `orm:"size(200)"`  // 矿机算力
	Power    string      `orm:"size(200)"`  // 矿机功率
	Eer      string      `orm:"size(200)"`  // 矿机能效比
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

func GetAllkj() (int64, []Goods, error) {
	o := orm.NewOrm()
	var info []Goods
	num, err := o.Raw("SELECT * FROM goods ").QueryRows(&info)
	return num, info, err
}

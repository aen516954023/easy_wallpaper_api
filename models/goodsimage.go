package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(GoodsImage))
}

// 商品图片表
type GoodsImage struct {
	Id       int
	Image    string    // 商品图片
	GoodsSKU *GoodsSKU `orm:"rel(fk)"` // 商品SKU
}

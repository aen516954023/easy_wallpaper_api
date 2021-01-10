package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(IndexTypeGoodsBanner))
}

// 首页分类商品展示表
type IndexTypeGoodsBanner struct {
	Id          int
	GoodsType   *GoodsType `orm:"rel(fk)"`
	GoodsSKU    *GoodsSKU  `orm:"rel(fk)"`
	DisplayType int        `orm:"default(1)"` // 展示类型：0 代表文字；1 代表图片
	Index       int        `orm:"default(0)"` // 展示顺序
}

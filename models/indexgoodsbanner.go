package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(IndexGoodsBanner))
}

// 首页轮播商品展示表
type IndexGoodsBanner struct {
	Id       int
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
	Image    string    // 商品图片
	Index    int       `orm:"default(0)"` // 商品展示顺序

}

type IndexBanner struct {
	Id         int
	GoodsSKUId int
	Image      string // 商品图片
	Index      int    // 商品展示顺序
}

// 获取首页轮播数据
func GetBannerAll() (int64, []IndexBanner, error) {
	var data []IndexBanner
	o := orm.NewOrm()
	num, err := o.Raw("SELECT * FROM index_goods_banner").QueryRows(&data)
	return num, data, err
}

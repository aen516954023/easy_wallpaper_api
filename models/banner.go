package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(EBanner))
}

// 首页轮播商品展示表
type EBanner struct {
	Id         int
	BannerName string
	Images     string // 图片
	Sort       int    // 展示顺序
}

// 获取首页轮播数据
func GetBannerAll() (int64, []EBanner, error) {
	var data []EBanner
	o := orm.NewOrm()
	num, err := o.Raw("SELECT * FROM e_banner").QueryRows(&data)
	return num, data, err
}

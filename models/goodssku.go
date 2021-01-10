package models

import (
	"easy_wallpaper_api/util"
	"fmt"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(GoodsSKU))
}

// 商品SKU表
type GoodsSKU struct {
	Id                   int
	Goods                *Goods                  `json:"goods" orm:"rel(fk)"`       // 商品SPU
	GoodsType            *GoodsType              `json:"goods_type" orm:"rel(fk)"`  // 商品类型
	GoodsCycle           *GoodsCycle             `json:"goods_cycle" orm:"rel(fk)"` // 商品套餐类型
	GoodsMode            *GoodsMode              `json:"goods_mode" orm:"rel(fk)"`  // 商品套餐类型
	Name                 string                  `orm:"size(50)"`                   // 商品名称
	Desc                 string                  `orm:"size(100)"`                  // 商品简介
	Price                float64                 // 商品价格
	Fee                  float64                 // 电费
	Unite                int                     `orm:"size(20)"` // 算力数量
	Image                string                  // 商品图片
	Stock                int                     `orm:"default(1)"`    // 商品库存
	Sales                int                     `orm:"default(0)"`    // 商品销量
	StaticIncome         float64                 `orm:"default(0.00)"` // 静态收益
	Status               int                     `orm:"default(1)"`    // 商品状态：是否有效，默认有效
	BeginTime            string                  `orm:"default(0)"`    // 预计生效时间
	Time                 string                  `orm:"auto_now_add"`  // 添加时间
	GoodsImage           []*GoodsImage           `orm:"reverse(many)"` // 商品图片
	IndexGoodsBanner     []*IndexGoodsBanner     `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
	OrderGoods           []*OrderGoods           `orm:"reverse(many)"`
}

func Get_goods_by_type(id int) []GoodsSKU {
	o := orm.NewOrm()
	var goodssku []GoodsSKU
	_, error := o.QueryTable("GoodsSKU").Filter("GoodsType__Id", id).RelatedSel("Goods").RelatedSel("GoodsType").RelatedSel("GoodsMode").RelatedSel("GoodsCycle").All(&goodssku)
	if error != nil {
		fmt.Printf("err", error)
	}
	return goodssku
}

func GetGoodsListByTj(where map[string]string) []GoodsSKU {
	o := orm.NewOrm()
	var goodssku []GoodsSKU
	var p orm.QuerySeter = o.QueryTable("GoodsSKU").RelatedSel("Goods").RelatedSel("GoodsType").RelatedSel("GoodsMode").RelatedSel("GoodsCycle")
	p = GetGoodsListWhere(p, where)
	p.All(&goodssku)
	//fmt.Println("goodssku", goodssku)
	return goodssku
}

//拼接where条件
func GetGoodsListWhere(p orm.QuerySeter, where map[string]string) orm.QuerySeter {
	if v, ok := where["goods_cycle_id"]; ok && v != "" {
		goods_cycle_id := util.TrimString(v)
		p = p.Filter("goodscycle__id", goods_cycle_id)
	}
	if v, ok := where["goods_type_id"]; ok && v != "" {
		goods_type_id := util.TrimString(v)
		p = p.Filter("goodstype__id", goods_type_id)
	}

	if v, ok := where["goods_mode_id"]; ok && v != "" {
		goods_mode_id := util.TrimString(v)
		p = p.Filter("GoodsMode__Id", goods_mode_id)
	}
	if v, ok := where["goods_id"]; ok && v != "" {
		goods_id := util.TrimString(v)
		p = p.Filter("Goods__Id", goods_id)
	}
	return p
}
func Get_goods_by_goodsskuid(id int) []GoodsSKU {
	o := orm.NewOrm()
	var goodssku []GoodsSKU
	_, error := o.QueryTable("GoodsSKU").Filter("Id", id).RelatedSel("Goods").RelatedSel("GoodsType").RelatedSel("GoodsMode").RelatedSel("GoodsCycle").All(&goodssku)
	if error != nil {
		fmt.Printf("err", error)
	}
	return goodssku
}

// 首页商品数据
func GetRecommendList(limit int) (int64, []GoodsSKU, error) {
	o := orm.NewOrm()
	var data []GoodsSKU
	num, err := o.QueryTable("GoodsSKU").RelatedSel("Goods").RelatedSel("GoodsType").RelatedSel("GoodsMode").RelatedSel("GoodsCycle").Limit(limit).All(&data)
	return num, data, err
}

// 获取单个sku产品详情
func GetGoodSKUDataOne(id int) (GoodsSKU, error) {
	o := orm.NewOrm()
	var data GoodsSKU
	err := o.QueryTable("GoodsSKU").RelatedSel("GoodsCycle").Filter("id", id).One(&data)
	return data, err
}

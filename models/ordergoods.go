package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(OrderGoods))
}

// 订单表
// 订单商品表
type OrderGoods struct {
	Id        int
	OrderInfo *OrderInfo `orm:"rel(fk)"`
	GoodsSKU  *GoodsSKU  `orm:"rel(fk)"`
	Count     int        `orm:"default(1)"` // 商品数量
	Price     float64    // 商品价格
	Comment   string     `orm:"default('');size(200)"` // 评论内容
}

func GetOrderInfo(id int) (OrderGoods, error) {
	o := orm.NewOrm()
	var data OrderGoods
	err := o.QueryTable("order_goods").
		RelatedSel("GoodsSKU").
		RelatedSel("OrderInfo").
		Filter("OrderInfo__Id", id).
		Filter("OrderInfo__OrderStatus", 1).
		One(&data)
	return data, err
}

type PayParams struct {
	OrderId      string //订单号
	TotalPrice   float64
	TransitPrice float64
	Id           int
	Name         string
	Price        float64
}

// 支付订单查询
func GetPayOrderInfo(sn string, id int) (PayParams, error) {
	o := orm.NewOrm()
	var data PayParams
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select("i.order_id,i.total_price,i.transit_price,sku.id,sku.name,sku.price").
		From("order_goods g").
		LeftJoin("order_info i").On("i.id = g.order_info_id").
		LeftJoin("user u").On("i.user_id = u.id").
		LeftJoin("goods_s_k_u sku").On("sku.id = g.goods_s_k_u_id").
		//Where("oi.order_id = ?")
		Where("i.order_id = ? and u.id = ?")
	// 导出 SQL 语句
	sql := qb.String()
	//fmt.Println(sql)
	// 执行 SQL 语句
	err := o.Raw(sql, sn, id).QueryRow(&data)
	return data, err
}

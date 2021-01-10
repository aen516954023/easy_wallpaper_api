package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(OrderInfo))
}

// 订单表
type OrderInfo struct {
	Id             int
	OrderId        string        `orm:"size(20);unique"` // 订单号
	User           *User         `orm:"rel(fk)"`
	Address        *Address      `orm:"rel(fk)"` // 收货地址
	PayMethod      int           // 支付方式
	TotalCount     int           `orm:"default(1)"` // 商品数量
	TotalPrice     float64       // 商品总价（包含运费）
	TransitPrice   float64       `orm:"default(0.00)"` // 运费
	CycleDay       int           // 电费天数
	Unite          int           // 算力值
	OrderStatus    int           `orm:"default(1)"`           // 支付状态: 1未支付， 2支付成功
	TradeNo        string        `orm:"size(20);default('')"` // 支付编号
	CreateAt       string        // 时间
	CreateTime     int64         // 创建时间
	ExpirationTime int64         // 创建时间
	GoodsSkuName   string        // sku名称
	OrderGoods     []*OrderGoods `orm:"reverse(many)"`
}

// 获取用户订单列表
func GetAllOrders(uid int) (int64, []OrderInfo, error) {
	o := orm.NewOrm()
	var data []OrderInfo
	num, err := o.QueryTable("order_info").Filter("user_id", uid).All(&data)
	return num, data, err
}

// 获取未支付订单列表
func GetUnpaidOrders() (int64, []OrderInfo, error) {
	o := orm.NewOrm()
	var data []OrderInfo
	num, err := o.QueryTable("order_info").Filter("order_status", 1).All(&data)
	return num, data, err
}

//获取符合方法收益条件的所有订单
func GetAllRorders(currenttime int64) []OrderInfo {
	var orders []OrderInfo
	begintime := currenttime - 60*60*24
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	qb.Select("*").
		From("order_info").
		Where("create_time < ? and order_status = ? and expiration_time > ? ").
		OrderBy("id").Desc()
	// 导出 SQL 语句
	sql := qb.String()
	// 执行 SQL 语句
	o := orm.NewOrm()
	o.Raw(sql, begintime, 2, currenttime).QueryRows(&orders)
	return orders
}

//获取所有订单的数量
func GetOrderAllNum(uid int) (int64, error) {
	o := orm.NewOrm()
	var orderinfo []OrderInfo
	num, err := o.Raw("SELECT id FROM order_info WHERE user_id = ?", uid).QueryRows(&orderinfo)
	if err == nil {
		return num, err
	} else {
		return 0, err
	}

}

//获取有效订单的数量
func GetOrderEffectiveNum(uid int, currenttime int64) (int64, int, error) {
	o := orm.NewOrm()
	begintime := currenttime - 60*60*24
	var orderinfo []OrderInfo
	num, err := o.Raw("SELECT total_count,unite FROM order_info WHERE user_id = ? and create_time < ? and order_status = ? and expiration_time > ?", uid, begintime, 2, currenttime).QueryRows(&orderinfo)
	if err == nil {
		//获取总算力
		allunite := 0
		for _, value := range orderinfo {
			allunite += value.TotalCount * value.Unite
		}
		return num, allunite, err
	} else {
		return 0, 0, err
	}

}

// 回调中查询订单信息
func GetNotifyOrderInfo(trade_no string) (int64, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("order_info").Filter("trade_no", trade_no).Count()
	return num, err
}

// 更新支付单号
func ModifyOrderTradeNo(order_sn, trade_no string) (bool, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("order_info").
		Filter("order_id", order_sn).
		Filter("order_status", 1).
		Update(orm.Params{
			"trade_no": trade_no,
		})
	if err == nil && num > 0 {
		return true, err
	}
	return false, err

}

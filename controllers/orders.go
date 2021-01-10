package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type Orders struct {
	Base
}

// @Title 订单详情
// @Description 订单详情
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_order_info [post]
func (this *Orders) Index() {
	orderId, err := this.GetInt("order_id") // 订单id
	// 参数效验
	if err != nil {
		this.Data["json"] = ReturnError(40001, "id参数不合法")
		this.ServeJSON()
		this.StopRun()
	}

	info, errs := models.GetOrderInfo(orderId)
	if errs != nil {
		this.Data["json"] = ReturnError(40001, "订单不存在")
		this.ServeJSON()
		this.StopRun()
	}

	infoMap := make(map[string]interface{})
	infoMap["title"] = info.GoodsSKU.Name
	infoMap["goods_total_price"] = info.OrderInfo.TotalPrice
	infoMap["price"] = info.GoodsSKU.Price
	infoMap["num"] = info.OrderInfo.TotalCount
	infoMap["power_fee"] = info.GoodsSKU.Fee
	infoMap["total_fee"] = info.OrderInfo.TransitPrice
	infoMap["day"] = info.OrderInfo.CycleDay
	infoMap["order_sn"] = info.OrderInfo.OrderId
	infoMap["time"] = info.OrderInfo.CreateAt

	this.Data["json"] = ReturnSuccess(0, "success", infoMap, 1)
	this.ServeJSON()
}

// @Title 提交订单
// @Description 提交订单
// @Param	goods_id		query 	int	true		"the goods sku id"
// @Param	nums		query 	int	true		"the buy num"
// @Param	total_price		query 	float64	true		"the total_price"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_order [post]
func (this *Orders) SaveOrder() {

	// 获取参数
	goodsId, err := this.GetInt("goods_id")               // 商品id
	nums, num_err := this.GetInt("nums")                  // 购买数量
	totalPrice, price_err := this.GetFloat("total_price") // 订单总金额

	// 参数效验
	if err != nil {
		this.Data["json"] = ReturnError(40001, "参数不合法")
		this.ServeJSON()
		this.StopRun()
	}

	if nums <= 0 || num_err != nil {
		this.Data["json"] = ReturnError(40001, "购买数量不合法")
		this.ServeJSON()
		this.StopRun()
	}

	if price_err != nil {
		this.Data["json"] = ReturnError(40001, "金额参数不合法")
		this.ServeJSON()
		this.StopRun()
	}
	// 通过商品id 查询相关数据，及库存 及相关产品信息
	info, _ := models.GetGoodSKUDataOne(goodsId)

	//============= 逻辑效验 ==================
	//1. 产品是否下架
	if info.Status != 1 {
		this.Data["json"] = ReturnError(40001, "该产品已下架，请选择其它产品")
		this.ServeJSON()
		this.StopRun()
	}
	// 2.库存判断 ,处理超卖 Todo
	if info.Sales >= info.Stock {
		this.Data["json"] = ReturnError(40007, "已售磬")
		this.ServeJSON()
		this.StopRun()
	} else if nums > (info.Stock - info.Sales) {
		this.Data["json"] = ReturnError(40006, "库存不足")
		this.ServeJSON()
		this.StopRun()
	}
	// 3. 判断params中的总金额 与 产品对应的金额 是否相符

	// 4. ...

	// 并发处理 Todo

	//开启事务
	o := orm.NewOrm()
	beginError := o.Begin()
	if beginError != nil {
		this.Data["json"] = ReturnError(40000, "事务异常")
		this.ServeJSON()
		return
	}
	var user models.User
	user.Id = this.CurrentLoginUser.Id // 获取用户uid
	var address models.Address
	address.Id = 0 // Todo 通过uid及币种信息获取用户对应提币币种地址id
	var orderInfo models.OrderInfo
	orderInfo.OrderId = CreateRandOrderOn() // 生成订单号
	orderInfo.TotalCount = nums
	orderInfo.TotalPrice = totalPrice
	orderInfo.TransitPrice = info.Fee * float64(info.GoodsCycle.Day)
	orderInfo.CycleDay = info.GoodsCycle.Day
	orderInfo.Unite = info.Unite
	orderInfo.GoodsSkuName = info.Name
	orderInfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	orderInfo.CreateTime = strToUnixTime(time.Now().Format("2006-01-02 15:04:05"))
	orderInfo.ExpirationTime = strToUnixTime(time.Now().Format("2006-01-02 15:04:05")) + 60*60*24*int64(info.GoodsCycle.Day)
	orderInfo.OrderStatus = 1
	orderInfo.Address = &address
	orderInfo.User = &user
	id, errors := o.Insert(&orderInfo)
	if errors != nil {
		rErr := o.Rollback()
		if rErr != nil {
			logs.Error("事务回滚出错")
		}
		logs.Error(errors)
		this.Data["json"] = ReturnError(40001, "订单提交错误")
		this.ServeJSON()
		return
	}
	orderInfo.Id = int(id)
	var goodsSKU models.GoodsSKU
	goodsSKU.Id = goodsId
	var orderGoods models.OrderGoods
	orderGoods.Count = nums
	orderGoods.Price = totalPrice
	orderGoods.OrderInfo = &orderInfo
	orderGoods.GoodsSKU = &goodsSKU

	_, errs := o.Insert(&orderGoods)
	if errs != nil {
		rErr := o.Rollback()
		if rErr != nil {
			logs.Error("事务回滚出错")
		}
		logs.Error(errs)
		this.Data["json"] = ReturnError(40001, "订单提交错误2")
		this.ServeJSON()
		return
	}
	// 订单提交成功 更新库存+nums
	goodsSKU.Sales = info.Sales + nums
	_, uErr := o.Update(&goodsSKU, "sales")
	if uErr != nil {
		rErr := o.Rollback()
		if rErr != nil {
			logs.Error("事务回滚出错")
		}
		logs.Error(uErr)
		this.Data["json"] = ReturnError(40001, "订单提交错误3")
		this.ServeJSON()
		return
	}
	ok := o.Commit()
	if ok != nil {
		logs.Error("事务提交出错")
	}
	this.Data["json"] = ReturnSuccess(0, "success", id, 1)
	this.ServeJSON()
}

// 取消订单

// @Title 支付接口
// @Description 支付详情
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /pay_order [post]
func (this *Orders) PayOrder() {
	orderId := this.GetString("order_id")
	payType, err := this.GetInt("pay_type")

	if err != nil {
		this.Data["json"] = ReturnError(40001, "支付单号不能为空")
		this.ServeJSON()
		this.StopRun()
	}
	// 查询订单数据
	info, errs := models.GetPayOrderInfo(orderId, this.CurrentLoginUser.Id)
	if errs != nil {
		this.Data["json"] = ReturnError(40001, "订单错误或订单不存在")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(info)

	switch payType {
	case 1: // 信用卡
		postdata := make(map[string]interface{})
		postdata["orders_code"] = info.OrderId                              //订单号
		postdata["order_total"] = (info.TotalPrice + info.TransitPrice)     //支付总金额
		postdata["currency_code"] = "USD"                                   //币种，例：美金USD
		postdata["order_total_usd"] = (info.TotalPrice + info.TransitPrice) //总折算美金金额
		postdata["notify_url"] = Config("notify_url")                       //支付结果回调地址 http://localhost:8055/notify
		postdata["products_id"] = info.Id                                   //产品id
		postdata["products_name"] = info.Name                               //产品名称
		postdata["products_price"] = info.Price                             //产品价格
		postdata["products_price_usd"] = info.Price                         //产品折算美金价格
		data := GetOrderUrl(postdata)

		// 更新支付单号
		boolVal, errVal := models.ModifyOrderTradeNo(info.OrderId, data["orders_id"].(string))
		fmt.Println(boolVal)
		if errVal == nil && boolVal {
			this.Data["json"] = ReturnSuccess(0, "success", data, 1)
			this.ServeJSON()
		} else {
			logs.Error("支付请求错误:" + fmt.Sprintf("%s", errVal))
			this.Data["json"] = ReturnError(40003, "支付请求错误")
			this.ServeJSON()
		}
		break
	default:
		this.Data["json"] = ReturnError(40004, "支付通道暂未开通")
		this.ServeJSON()
	}

}

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

// @Title 订单大厅
// @Description 订单大厅，只有师傅用户能够查看
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_orders_all [post]
func (this *Orders) Index() {

	// 获取地理位置 查询当前地区订单信息 Todo

	num, data, err := models.GetOrdersAll()
	if err == nil && num > 0 {
		this.Data["json"] = ReturnSuccess(0, "success", data, num)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 提交订单
// @Description 普通用户提交订单
// @Param	token		header 	string	true		"the token"
// @Param	address		query 	string	true		"the address"
// @Param	construction_time		query 	string	true	"the construction_time"
// @Param	types		query 	int	true		"the types"
// @Param	is_materials		query 	int	true		"the is_materials"
// @Param	area		query 	float64	true		"the area"
// @Param	is_tear_of_old_wallpaper		query 	int	true		"the is_tear_of_old_wallpaper"
// @Param	basement_membrane		query 	int	true		"the basement_membrane"
// @Param	desc		query 	string	true		"the desc"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_order [post]
func (this *Orders) SaveOrder() {

	// 获取参数
	Address := this.GetString("address")                       // 施工地址
	constructionTime := this.GetString("construction_time")    // 施工时间
	types, typesErr := this.GetInt("types")                    // 服务类型
	materials, _ := this.GetInt("is_materials")                // 是否提供物料
	area, areaErr := this.GetFloat("area")                     // 面积
	oldWallpaper, _ := this.GetInt("is_tear_of_old_wallpaper") // 是否清除旧纸
	basementMembrane, _ := this.GetInt("basement_membrane")    // 是否刷基膜
	moreDescription := this.GetString("desc")                  // 需求描述
	// 图片上传

	// 用户参数
	userId := this.CurrentLoginUser.Id // 当前用户id
	//phone := "15938755991"
	orderType, orderTypesErr := this.GetInt("order_type") // 订单类型
	if orderTypesErr == nil && orderType == 2 {
		orderType = 2
	}

	// 订单类型为直接选择师傅
	workerId, workerErr := this.GetInt("worker_id") // 师傅id
	if workerErr != nil {
		workerId = 0
	}

	// 参数效验
	if Address == "" {
		this.Data["json"] = ReturnError(40001, "地址不能为空")
		this.ServeJSON()
		this.StopRun()
	}

	if constructionTime == "" {
		this.Data["json"] = ReturnError(40001, "请选择施工时间")
		this.ServeJSON()
		this.StopRun()
	}

	if typesErr != nil {
		this.Data["json"] = ReturnError(40001, "服务类型参数错误")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(types)

	if areaErr != nil || area < 1 {
		this.Data["json"] = ReturnError(40001, "请选择有效的施工面积")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(area)
	var orderinfo models.EOrders
	//orderinfo.OrderSn = CreateRandOrderOn() // 生成订单号
	orderinfo.MId = int(userId)
	orderinfo.WorkerId = workerId
	orderinfo.Area = area
	orderinfo.BasementMembrane = basementMembrane
	orderinfo.IsTearOfOldWallpaper = oldWallpaper
	orderinfo.IsMateriel = materials
	orderinfo.MoreDescription = moreDescription
	orderinfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	orderinfo.Status = 1
	orderinfo.OrderType = orderType
	orderinfo.Address = Address
	fmt.Println(orderinfo)
	//panic("orderinfo")
	o := orm.NewOrm()
	id, errors := o.Insert(&orderinfo)
	fmt.Println(errors)
	if errors == nil {
		this.Data["json"] = ReturnSuccess(0, "success", id, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40001, "订单提交失败")
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

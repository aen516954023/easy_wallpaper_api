package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"log"
)

/**
订单步骤业务逻辑
*/
type OrderStep struct {
	Base
}

// @Title 确认师傅
// @Description 用户确认选择师傅
// @Param	order_id		query 	int	true		"the order id"
// @Param	w_id		query 	int	true		"the worker id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_master_worker [post]
func (this *OrderStep) ConfirmMasterWorker() {
	// order_step 表新增选择的师傅记录
	// 更新师傅参与表中其它参与的师傅的状态
	orderId, _ := this.GetInt("order_id")
	workerId, _ := this.GetInt("w_id")
	if models.InsertOrdersStep(orderId, int(this.CurrentLoginUser.Id), workerId) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 0)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40003, "选择失败，请稍后再试")
		this.ServeJSON()
	}
}

// @Title 基础报价
// @Description 师傅发起基础报价
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /advance_order [post]
//
func (this *OrderStep) AdvanceOrder() {
	// 接收参数
	orderId, _ := this.GetInt("order_id")
	mId := int(this.CurrentLoginUser.Id)
	serviceType, _ := this.GetInt("service_type")
	constructionType, _ := this.GetInt("construction_type")
	price, _ := this.GetFloat("price")
	unit, _ := this.GetInt("unit")
	info := this.GetString("info")
	depositPrice, _ := this.GetFloat("deposit_price")
	// 参数效验 Todo
	fmt.Println(orderId)
	valid := validation.Validation{}
	valid.Required(orderId, "order_id")
	valid.Required(serviceType, "serviceType")
	valid.Required(constructionType, "constructionType")
	valid.Required(price, "price")
	valid.Required(info, "info")
	valid.Required(depositPrice, "depositPrice")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			//log.Fatal(err.Key, err.Message)
			this.Data["json"] = ReturnError(40000, err.Key+err.Message)
			this.ServeJSON()
			return
		}
	}
	//查询当前用户的师傅id
	workerInfo, workerInfoErr := models.GetMasterWorkerInfo(this.CurrentLoginUser.Id)
	if workerInfoErr == nil && workerInfo.Id > 0 {
		// 新增订单步骤 并更新订单表对应的相关信息
		if models.InsertOrdersStepTwo(orderId, mId, workerInfo.Id, serviceType, constructionType, unit, price, depositPrice, info) {
			this.Data["json"] = ReturnSuccess(0, "success", "", 0)
			this.ServeJSON()
			return
		} else {
			this.Data["json"] = ReturnError(40003, "操作失败，请稍后再试")
			this.ServeJSON()
			return
		}
	} else {
		this.Data["json"] = ReturnError(40004, "操作失败，师傅信息不存在")
		this.ServeJSON()
	}
}

// @Title 确认基础报价
// @Description 用户确认基础报价，并生成预支付订单
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_advance_order [post]
func (this *OrderStep) ConfirmAdvanceOrder() {
	oId, _ := this.GetInt("order_id")
	//通过订单id 查询订单信息,与基础报价信息
	data, err := models.GetOrderOfStepInfo(oId, 1)
	if err != nil {
		this.Data["json"] = ReturnError(40000, "订单信息不存在")
		this.ServeJSON()
		return
	}
	//判断订单是否存在
	count, _ := models.GetOrderPayEmpty(oId, int(this.CurrentLoginUser.Id), 1)
	if count > 0 {
		this.Data["json"] = ReturnError(40001, "请务重复确认基础订单")
		this.ServeJSON()
		return
	}
	//生成支付订单信息  如果支付成功 在回调方法中添加一条订单步骤 状态为2
	insertId, retValErr := models.InsertOrderPayInfo(CreateRandOrderOn(), oId, int(this.CurrentLoginUser.Id), 1, data.DepositPrice)
	if retValErr == nil && insertId > 0 {
		// 更新当前步骤支付状态与支付订单id
		boolVal, _ := models.UpdateOrderStepPayStatus(oId, this.CurrentLoginUser.Id, 1, 1, insertId)
		if boolVal {
			this.Data["json"] = ReturnSuccess(0, "success", insertId, 1)
			this.ServeJSON()
			return
		}
		this.Data["json"] = ReturnError(40003, "订单已确认，未支付")
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40004, "操作失败，请稍后再试")
	this.ServeJSON()
}

// @Title 实际报价
// @Description 师傅现场量房，发起实际报价
// @Param	order_id		query 	int	true		"the order id 订单id"
// @Param	w_id		query 	int	true		"the worker id 师傅id"
// @Param	service_type		query 	int	true		"the service_type 服务类型"
// @Param	construction_type		query 	int	true		"the construction_type 施工类型"
// @Param	price		query 	float64	true		"the price 报价"
// @Param	unit		query 	int	true		"the unit类型"
// @Param	info		query 	string	true		"the info"
// @Param	deposit_price		query 	float	true		"the deposit_price 定金"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /actual_offer [post]
func (this *OrderStep) ActualOffer() {
	// 接收参数
	orderId, _ := this.GetInt("order_id")
	mId := int(this.CurrentLoginUser.Id)
	wId, _ := this.GetInt("w_id")
	serviceType, _ := this.GetInt("service_type")
	constructionType, _ := this.GetInt("construction_type")
	price, _ := this.GetFloat("price")
	unit, _ := this.GetInt("unit")
	info := this.GetString("info")
	depositPrice, _ := this.GetFloat("deposit_price")
	// 参数效验 Todo
	valid := validation.Validation{}
	valid.Required(orderId, "order_id")
	valid.Required(wId, "w_id")
	valid.Required(serviceType, "serviceType")
	valid.Required(constructionType, "constructionType")
	valid.Required(price, "price")
	valid.Required(info, "info")
	valid.Required(depositPrice, "depositPrice")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			log.Fatal(err.Key, err.Message)
			this.Data["json"] = ReturnError(40000, err.Message)
			this.ServeJSON()
			return
		}
	}
	// 新增订单步骤 并更新订单表对应的相关信息
	if models.InsertOrdersStepThree(orderId, mId, wId, serviceType, constructionType, unit, price, depositPrice, info) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 0)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40003, "操作失败，请稍后再试")
		this.ServeJSON()
	}
}

// @Title 实际价格确认
// @Description 用户确认实际价格，并生成支付订单
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_actual_offer [post]
func (this *OrderStep) ConfirmActualOffer() {

	oId, _ := this.GetInt("order_id")
	//通过订单id 查询订单信息,实际报价信息
	data, err := models.GetOrderOfStepInfo(oId, 2)
	if err != nil {
		this.Data["json"] = ReturnError(40000, "订单信息不存在")
		this.ServeJSON()
		return
	}
	//生成支付订单信息
	insertId, retValErr := models.InsertOrderPayInfo(CreateRandOrderOn(), oId, int(this.CurrentLoginUser.Id), 2, data.TotalPrice)
	if retValErr == nil && insertId > 0 {
		this.Data["json"] = ReturnSuccess(0, "success", insertId, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40004, "操作失败，请稍后再试")
	this.ServeJSON()
}

// @Title 验收
// @Description 施工完成，师傅发起验收通知
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /Acceptance [post]
func (this *OrderStep) Acceptance() {
	oId, _ := this.GetInt("order_id")
	//通过订单id 查询订单信息,实际报价信息
	data, err := models.GetOrderOfStepInfo(oId, 4)
	if err != nil {
		this.Data["json"] = ReturnError(40000, "订单信息不存在")
		this.ServeJSON()
		return
	}
	if models.ModifyStepStatus(data.OId, int(this.CurrentLoginUser.Id), 5) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40003, "操作失败，请稍后再试")
	this.ServeJSON()
}

// @Title 验收确认
// @Description 用户验收确认
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /Acceptance [post]
func (this *OrderStep) ConfirmAcceptance() {
	oId, _ := this.GetInt("order_id")
	//通过订单id 查询订单信息,实际报价信息
	data, err := models.GetOrderOfStepInfo(oId, 5)
	if err != nil {
		this.Data["json"] = ReturnError(40000, "订单信息不存在")
		this.ServeJSON()
		return
	}
	if models.ModifyStepStatus(data.OId, int(this.CurrentLoginUser.Id), 6) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40003, "操作失败，请稍后再试")
	this.ServeJSON()
}

// @Title 支付页接口
// @Description 支付页接口数据
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_id		query 	int	true		"the pay id"
// @Param	flag		query 	int	true		"the flag 0|1"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_pay_info [post]
func (this *OrderStep) GetPayInfo() {
	oId, _ := this.GetInt("order_id")
	pId, _ := this.GetInt("pay_id")
	flag, _ := this.GetInt("flag")
	if oId == 0 || pId == 0 {
		this.Data["json"] = ReturnError(40001, "参数错误")
		this.ServeJSON()
		return
	}

	data := make(map[string]interface{})
	payData, payErr := models.GetOrdersPayInfo(pId)
	if payErr == nil && payData.Id > 0 {
		data["order_sn"] = payData.OrderSn
		data["total_price"] = payData.TotalPrice
	}
	status := 1
	if flag > 0 {
		status = 3
	}
	orderData, orderErr := models.GetOrderOfStepInfo(oId, status)
	if orderErr == nil && orderData.Id > 0 {
		data["service_type"] = orderData.ServiceType
		data["construction_type"] = orderData.ConstructionType
		data["area"] = orderData.Area
		data["price"] = orderData.Price
		data["deposit_price"] = orderData.DepositPrice
		data["unit"] = orderData.Unit
		data["address"] = orderData.OId
	}

	this.Data["json"] = ReturnSuccess(0, "success", data, 1)
	this.ServeJSON()
}

package controllers

import (
	"easy_wallpaper_api/models"
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
// @router /advance_order [post]
//
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
			this.Data["json"] = ReturnError(40000, err.Key+err.Message)
			this.ServeJSON()
			return
		}
	}
	// 新增订单步骤 并更新订单表对应的相关信息
	if models.InsertOrdersStepTwo(orderId, mId, wId, serviceType, constructionType, unit, price, depositPrice, info) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 0)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40003, "操作失败，请稍后再试")
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
	//生成支付订单信息
	insertId, retValErr := models.InsertOrderPayInfo(CreateRandOrderOn(), oId, int(this.CurrentLoginUser.Id), 1, data.DepositPrice)
	if retValErr == nil && insertId > 0 {
		this.Data["json"] = ReturnSuccess(0, "success", insertId, 1)
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

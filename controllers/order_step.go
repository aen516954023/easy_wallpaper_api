package controllers

/**
订单步骤业务逻辑
*/
type OrderStep struct {
	Base
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
	// 2 更新订单状态
}

// @Title 确认基础报价
// @Description 用户确认基础报价，并生成预支付订单
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_advance_order [post]
//
func (this *OrderStep) ConfirmAdvanceOrder() {
	// 1. 生成支付订单信息

	// 2. 回调中处理逻辑
}

// @Title 实际报价
// @Description 师傅现场量房，发起实际报价
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /actual_offer [post]
func (this *OrderStep) ActualOffer() {
	// 1 接收参数，要修改的订单id, 修改后的订单参数
	// 2 更新订单信息表 相关需求变更参数，更新订单状态
	// 3 新增全额支付订单信息表 支付总金额
	// 4 新增订单步骤表信息
}

// @Title 实际价格确认
// @Description 用户确认实际价格，并生成支付订单
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_actual_offer [post]
func (this *OrderStep) ConfirmActualOffer() {
	// 1 接收参数，要修改的订单id, 修改后的订单参数
	// 2 更新订单信息表 相关需求变更参数，更新订单状态
	// 3 新增全额支付订单信息表 支付总金额
	// 4 新增订单步骤表信息
}

// @Title 验收
// @Description 施工完成，师傅发起验收通知
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /Acceptance [post]
func (this *OrderStep) Acceptance() {

}

// @Title 验收确认
// @Description 用户验收确认
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /Acceptance [post]
func (this *OrderStep) ConfirmAcceptance() {

}

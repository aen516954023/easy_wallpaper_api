package controllers

import "easy_wallpaper_api/models"

type Workers struct {
	Base
}

// @Title 师傅接单
// @Description 师傅接单接口
// @Param	token		header 	string	true		"the token"
// @Param	order_id		query 	int	true		"the order id"
// @Param	m_id		query 	int	true		"the member id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /order_taking [post]
func (this *Workers) OrderTaking() {
	oId, _ := this.GetInt("order_id")
	if oId == 0 {
		this.Data["jsong"] = ReturnError(40001, "订单参数错误或不能为空")
		this.ServeJSON()
		return
	}
	mId, _ := this.GetInt("m_id")
	if mId == 0 {
		this.Data["jsong"] = ReturnError(40001, "会员参数错误或不能为空")
		this.ServeJSON()
		return
	}

	if models.InsertOrderTaking(oId, mId, int(this.CurrentLoginUser.Id)) {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	this.Data["jsong"] = ReturnError(40004, "接单失败，请稍后再试")
	this.ServeJSON()
}

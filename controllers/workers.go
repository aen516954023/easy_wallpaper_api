package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

type Workers struct {
	Base
}

// @Title 师傅入驻页
// @Description 师傅入驻申请页数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_settle_in_page [post]
func (this *Workers) SettleInPage() {
	//查看当前用户是否申请师傅，如果申请显示师傅状态
	data := make(map[string]int)
	data["is_master_worker"] = 0
	data["status"] = 0
	val, err := models.GetMasterWorkerInfo(this.CurrentLoginUser.Id)
	if err == nil && val.Id > 0 {
		data["is_master_worker"] = 1
		data["status"] = val.Status
	}
	this.Data["json"] = ReturnSuccess(0, "success", data, 1)
	this.ServeJSON()
}

// @Title 师傅入驻申请
// @Description 师傅入驻申请接口
// @Param	token		header 	string	true		"the token"
// @Param	username		query 	string	true		"the username"
// @Param	gender		query 	string	true		"the gender"
// @Param	mobile		query 	string	true		"the mobile"
// @Param	avatar		query 	string	true		"the avatar img"
// @Param	city		query 	string	true		"the service city"
// @Param	address		query 	string	true		"the address"
// @Param	exp		query 	int	true		"the exp"
// @Param	desc		query 	string	true		"the desc"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /apply [post]
func (this *Workers) Apply() {
	// 接收参数
	userName := this.GetString("username")
	gender, _ := this.GetInt("gender")
	mobile := this.GetString("mobile")
	city := this.GetString("city")
	address := this.GetString("address")
	exp, _ := this.GetInt("exp")
	desc := this.GetString("desc")
	avatar := this.GetString("image")
	// 参数校验
	if userName == "" {
		this.Data["json"] = ReturnError(40001, "用户名不能为空")
		this.ServeJSON()
		return
	}
	if gender == 0 {
		this.Data["json"] = ReturnError(40001, "性别不能为空")
		this.ServeJSON()
		return
	}
	if mobile == "" {
		this.Data["json"] = ReturnError(40001, "手机号码不能为空")
		this.ServeJSON()
		return
	}
	if city == "" {
		this.Data["json"] = ReturnError(40001, "服务城市不能为空")
		this.ServeJSON()
		return
	}
	if address == "" {
		this.Data["json"] = ReturnError(40001, "联系地址不能为空")
		this.ServeJSON()
		return
	}
	if exp == 0 {
		this.Data["json"] = ReturnError(40001, "施工经验不能为空")
		this.ServeJSON()
		return
	}
	if desc == "" {
		this.Data["json"] = ReturnError(40001, "个人描述不能为空")
		this.ServeJSON()
		return
	}

	// 写入数据
	boolVal, err := models.ApplyMasterWorker(this.CurrentLoginUser.Id, gender, exp, userName, mobile, city, address, desc, avatar, UnixTimeToSTr(time.Now().Unix()))
	if boolVal && err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	logs.Error("师傅入驻申请错误：" + fmt.Sprintf("%v", err))
	this.Data["json"] = ReturnError(40002, "提交审核失败,请稍后再试")
	this.ServeJSON()
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

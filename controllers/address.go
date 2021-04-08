package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
)

type Address struct {
	Base
}

// @Title 我的地址列表
// @Description 我的地址列表接口
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /address_list [post]
func (a *Address) Index() {
	mId := a.CurrentLoginUser.Id

	num, data, err := models.GetAddressList(mId)
	if err == nil && num > 0 {
		a.Data["json"] = ReturnSuccess(0, "success", data, num)
		a.ServeJSON()
	} else {
		a.Data["json"] = ReturnError(40000, "暂无记录")
		a.ServeJSON()
	}
}

// @Title 添加修改地址页
// @Description 添加修改地址页初始化页面数据接口
// @Param	a_id	query 	int 	true		"地址id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /address_page [post]
func (a *Address) AddressPage() {

	aId, _ := a.GetInt("a_id")
	// 修改操作
	if aId > 0 {
		info, err := models.GetAddressId(aId)
		if err == nil && info.Id > 0 {
			a.Data["json"] = ReturnSuccess(0, "success", info, 1)
			a.ServeJSON()
		} else {
			a.Data["json"] = ReturnError(40000, "暂无记录")
			a.ServeJSON()
		}
	} else {
		// 添加操作,只返回用户手机号码参数
		returnVal := make(map[string]string)
		returnVal["Phone"] = a.CurrentLoginUser.Phone
		a.Data["json"] = ReturnSuccess(0, "success", returnVal, 0)
		a.ServeJSON()
	}
}

// @Title 添加修改地址操作
// @Description 添加修改地址操作接口
// @Param	flag	query 	int 	true		"操作 1 添加 2修改"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /address_action [post]
func (a *Address) SaveAddress() {
	// 最多添加5个地址
	num, _, _ := models.GetAddressList(a.CurrentLoginUser.Id)
	if num >= 5 {
		a.Data["json"] = ReturnError(40005, "最多可添加5个地址")
		a.ServeJSON()
		return
	}
	// 接收参数 数据校验
	pm := a.Ctx.Request.Form
	valid := validation.Validation{}
	valid.Required(pm.Get("username"), "联系人不能为空")
	valid.Required(pm.Get("phone"), "地址方式不能为空")
	valid.Required(pm.Get("address_name"), "地址不能为空")
	valid.Required(pm.Get("address"), "地址详情不能为空")
	valid.Required(pm.Get("province"), "省份不能为空")
	valid.Required(pm.Get("city"), "城市不能为空")
	valid.Required(pm.Get("district"), "市区不能为空")
	valid.Required(pm.Get("house_number"), "门牌号码不能为空")

	latitude, _ := a.GetFloat("latitude")
	longitude, _ := a.GetFloat("longitude")
	isDefault, _ := a.GetInt("default")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			a.Data["json"] = ReturnError(40000, err.Key)
			a.ServeJSON()
			return
		}
	}
	// 处理逻辑
	flag, _ := a.GetInt("aid")
	if flag > 0 {
		//修改操作
		boolVal, err := models.ModifyAddress(
			flag,
			int(a.CurrentLoginUser.Id),
			isDefault,
			pm.Get("username"),
			pm.Get("phone"),
			pm.Get("address_name"),
			pm.Get("address"),
			pm.Get("province"),
			pm.Get("city"),
			pm.Get("district"),
			pm.Get("house_number"),
			latitude,
			longitude,
		)
		if err == nil && boolVal {
			a.Data["json"] = ReturnSuccess(0, "success", "", 1)
			a.ServeJSON()
		} else {
			logs.Error("地址修改失败原因:" + fmt.Sprintf("%v", err))
			a.Data["json"] = ReturnError(40003, "地址修改失败，请稍后再试")
			a.ServeJSON()
		}
	} else {
		//添加操作
		insertId, err := models.InsertAddress(
			int(a.CurrentLoginUser.Id),
			isDefault,
			pm.Get("username"),
			pm.Get("phone"),
			pm.Get("address_name"),
			pm.Get("address"),
			pm.Get("province"),
			pm.Get("city"),
			pm.Get("district"),
			pm.Get("house_number"),
			latitude,
			longitude,
		)
		if err == nil && insertId > 0 {
			a.Data["json"] = ReturnSuccess(0, "success", "", 1)
			a.ServeJSON()
		} else {
			logs.Error("地址添加失败原因:" + fmt.Sprintf("%v", err))
			a.Data["json"] = ReturnError(40003, "地址添加失败，请稍后再试")
			a.ServeJSON()
		}
	}
}

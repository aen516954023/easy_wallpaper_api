package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
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
	paramsMap := a.Ctx.Request.Form

	fmt.Println(paramsMap)

	//flag,_ :=
}

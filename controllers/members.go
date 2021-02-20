package controllers

import "easy_wallpaper_api/models"

type Members struct {
	Base
}

// @Title 我的
// @Description 个人中心页数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_members_center [post]
func (this *Members) Index() {
	data := make(map[string]interface{})
	//获取用户信息 昵称 头像 ip
	data["nickname"] = this.CurrentLoginUser.Nickname
	data["avatar_img"] = ""
	data["ip"] = 9293849
	// 师傅入驻接口相关参数
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

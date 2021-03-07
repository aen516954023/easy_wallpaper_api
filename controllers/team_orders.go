package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type TeamOrders struct {
	beego.Controller
}

// @Title 施工队官网表单收集
// @Description 施工队官网表单收集接口
// @Param	name		query 	string 	true		"the name"
// @Param	gender		query 	int 	true		"the gender"
// @Param	area		query 	string 	true		"the area"
// @Param	day		query 	int 	true		"the int"
// @Param	city		query 	string 	true		"the city"
// @Param	phone		query 	string 	true		"the phone"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /pc_data [post]
func (this *TeamOrders) Index() {
	name := this.GetString("name")
	gender, _ := this.GetInt("gender")
	area, _ := this.GetFloat("area")
	day, _ := this.GetInt("day")
	city := this.GetString("city")
	phone := this.GetString("phone")
	if name == "" {
		this.Data["json"] = ReturnError(40001, "称呼不能为空")
		this.ServeJSON()
		return
	}
	if gender == 0 {
		this.Data["json"] = ReturnError(40001, "性别不能为空")
		this.ServeJSON()
		return
	}
	if area <= 0 {
		this.Data["json"] = ReturnError(40001, "面积不能为空")
		this.ServeJSON()
		return
	}
	if day == 0 {
		this.Data["json"] = ReturnError(40001, "工期不能为空")
		this.ServeJSON()
		return
	}
	if city == "" {
		this.Data["json"] = ReturnError(40001, "城市不能为空")
		this.ServeJSON()
		return
	}
	if phone == "" {
		this.Data["json"] = ReturnError(40001, "手机号码不能为空")
		this.ServeJSON()
		return
	}

	boolVal, err := models.SaveTeamOrders(name, phone, city, gender, day, area)
	if !boolVal {
		logs.Error("获取报价失败：", fmt.Sprintf("%v", err))
		this.Data["json"] = ReturnError(40003, "获取失败，请稍后再试")
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnSuccess(0, "success", "", 1)
	this.ServeJSON()
}

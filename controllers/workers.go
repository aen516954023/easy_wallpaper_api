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

// @Title 师傅中心
// @Description 师傅中心数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_master_center [post]
func (this *Workers) GetMasterCenter() {
	// 查询当前用户师傅信息
	info, err := models.GetMasterWorkerInfo(this.CurrentLoginUser.Id)
	if err == nil && info.Id > 0 {
		returnVal := make(map[string]interface{})
		returnVal["username"] = info.Username
		returnVal["avatar"] = info.Image
		returnVal["is_real_name"] = info.IsRealName
		returnVal["is_exp"] = info.Exp
		returnVal["is_warranty"] = info.Warranty
		returnVal["order_count"] = 0        //可接单次数
		returnVal["account_money"] = "0.00" //可用余额

		this.Data["json"] = ReturnSuccess(0, "success", returnVal, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40001, "没有查询到相关数据")
	this.ServeJSON()
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

// @Title 师傅订单管理
// @Description 订单管理数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /order_list [post]
func (this *Workers) OrderList() {

	num, data, err := models.GetOrderMasterAll(this.CurrentLoginUser.Id)
	if err == nil && num > 0 {
		returnVal := make([]map[string]interface{}, num)
		for k := 0; k < len(returnVal); k++ {
			// 一定要加下面的nil判断  否则会报错 Handler crashed with error assignment to entry in nil map  map未赋值
			if returnVal[k] == nil {
				returnVal[k] = map[string]interface{}{}
			}
			returnVal[k]["Id"] = data[k].Id
			serviceTypeObj, _ := models.GetServiceType(int64(data[k].ServiceId))
			returnVal[k]["ServiceType"] = serviceTypeObj.TypeName
			var constructionData = []string{"主料+辅料+施工", "仅施工", "辅料+施工"}
			returnVal[k]["ConstructionType"] = constructionData[data[k].ConstructionType]
			returnVal[k]["Area"] = data[k].Area
			returnVal[k]["ConstructionTime"] = UnixTimeToSTr(int64(data[k].ConstructionTime))
			returnVal[k]["CreateAt"] = data[k].CreateAt
			returnVal[k]["OrderSn"] = data[k].OrderSn
		}
		this.Data["json"] = ReturnSuccess(0, "success", returnVal, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40000, "没有相关记录")
		this.ServeJSON()
	}
}

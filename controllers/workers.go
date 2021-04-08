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
		returnVal["id"] = info.Id
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
	// 接收订单状态参数
	status, _ := this.GetInt("status")
	s := 0
	if status != 0 {
		s = status
	}
	num, data, err := models.GetOrderMasterAll(this.CurrentLoginUser.Id, s)
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

// @Title 师傅信息修改页
// @Description 师傅信息修改页数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /edit_master_info [post]
func (w *Workers) EditMasterInfo() {

	data, err := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	if err == nil && data.Id > 0 {
		w.Data["json"] = ReturnSuccess(0, "success", data, 1)
		w.ServeJSON()
	} else {
		w.Data["json"] = ReturnError(40000, "未查询到记录")
		w.ServeJSON()
	}
}

// @Title 师傅信息修改提交
// @Description 师傅信息修改页提交接口
//// @Param	token		header 	string	true		"the token"
//// @Param	username		query 	string	true		"the username"
//// @Param	gender		query 	string	true		"the gender"
//// @Param	mobile		query 	string	true		"the mobile"
//// @Param	avatar		query 	string	true		"the avatar img"
//// @Param	city		query 	string	true		"the service city"
//// @Param	address		query 	string	true		"the address"
//// @Param	exp		query 	int	true		"the exp"
//// @Param	desc		query 	string	true		"the desc"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /save_edit_master_info [post]
func (this *Workers) SaveEditMasterInfo() {
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
	boolVal, err := models.SaveEditMasterInfo(this.CurrentLoginUser.Id, gender, exp, userName, mobile, city, address, desc, avatar, UnixTimeToSTr(time.Now().Unix()))
	if boolVal && err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	logs.Error("师傅资料修改错误：" + fmt.Sprintf("%v", err))
	this.Data["json"] = ReturnError(40002, "修改失败,请稍后再试")
	this.ServeJSON()
}

// @Title 案例发布
// @Description 发布新案例接口
//// @Param	token		header 	string	true		"the token"
//// @Param	image		query 	string	true		"the images"
//// @Param	desc		query 	string	true		"the desc"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /send_case [post]
func (w *Workers) SendCase() {
	desc := w.GetString("desc")
	image := w.GetString("image")
	if desc == "" {
		w.Data["json"] = ReturnError(40001, "描述不能为空")
		w.ServeJSON()
		return
	}
	if image == "" {
		w.Data["json"] = ReturnError(40001, "案例图不能为空")
		w.ServeJSON()
		return
	}
	info, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	if info.Id > 0 {
		boolVal, err := models.AddMasterCase(info.Id, desc, image)
		if err == nil && boolVal {
			w.Data["json"] = ReturnSuccess(0, "success", "", 1)
			w.ServeJSON()
			return
		}
		w.Data["json"] = ReturnError(40002, "发布失败，请稍后再试")
		w.ServeJSON()
		return
	}
	w.Data["json"] = ReturnError(40003, "发布失败,没有查询到相关师傅信息")
	w.ServeJSON()
}

// @Title 我的案例
// @Description 师傅中心-我的案例页面
//// @Param	token		header 	string	true		"the token"
//// @Param	page		query 	int 	false		"the page num"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /case_list [post]
func (w *Workers) CaseList() {
	page, _ := w.GetInt("page")
	if page == 0 {
		page = 1
	}
	pageNum := 10
	info, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	if info.Id > 0 {
		num, list, err := models.GetMasterCaseList(info.Id, page, pageNum)
		if err == nil {
			w.Data["json"] = ReturnSuccess(0, "success", list, num)
			w.ServeJSON()
		} else {
			w.Data["json"] = ReturnError(40000, "获取数据失败")
			w.ServeJSON()
		}
	} else {
		w.Data["json"] = ReturnError(40003, "师傅信息不存在")
		w.ServeJSON()
	}
}

// @Title 实名认证
// @Description 师傅中心-实名认证接口
//// @Param	token		header 	string	true		"the token"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /get_id_card [post]
func (w *Workers) GetIdCard() {
	info, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	if info.Id <= 0 {
		w.Data["json"] = ReturnError(40000, "师傅信息不存在")
		w.ServeJSON()
		return
	}
	idInfo, err := models.GetIdCard(info.Id)
	if err == nil && idInfo.Id > 0 {
		w.Data["json"] = ReturnSuccess(0, "success", idInfo, 1)
		w.ServeJSON()
	} else {
		w.Data["json"] = ReturnError(40001, "获取信息失败")
		w.ServeJSON()
	}
}

// @Title 实名认证提交
// @Description 师傅中心-实名认证提交接口
//// @Param	token		header 	string	true		"the token"
//// @Param	real_name	query 	string 	true		"the real_name"
//// @Param	id_card		query 	string 	true		"the id_card"
//// @Param	pos			query 	string 	true		"the image pos"
//// @Param	neg			query 	string 	true		"the image neg"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /save_id_card [post]
func (w *Workers) SaveIdCard() {
	realName := w.GetString("real_name")
	idCard := w.GetString("id_card")
	pos := w.GetString("pos")
	neg := w.GetString("neg")

	if realName == "" {
		w.Data["json"] = ReturnError(40000, "真实名字不能为空")
		w.ServeJSON()
		return
	}
	if idCard == "" {
		w.Data["json"] = ReturnError(40000, "身份证号码不能为空")
		w.ServeJSON()
		return
	}
	if pos == "" {
		w.Data["json"] = ReturnError(40000, "身份证正面照片不能为空")
		w.ServeJSON()
		return
	}
	if neg == "" {
		w.Data["json"] = ReturnError(40000, "身份证反面照片不能为空")
		w.ServeJSON()
		return
	}
	//获取当前用户师傅id
	workerInfo, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	//查询实名认证是否已提交
	cardInfo, _ := models.GetIdCard(workerInfo.Id)
	if cardInfo.Id > 0 {
		w.Data["json"] = ReturnError(40003, "实名认证已提交，勿重复提交")
		w.ServeJSON()
		return
	}
	boolVal, err := models.AddIdCard(workerInfo.Id, idCard, realName, pos, neg)
	if err == nil && boolVal {
		w.Data["json"] = ReturnSuccess(0, "success", "", 1)
		w.ServeJSON()
	} else {
		w.Data["json"] = ReturnError(40001, "提交失败，请稍后再试")
		w.ServeJSON()
	}
}

// @Title 经验认证
// @Description 师傅中心-经验认证接口
//// @Param	token		header 	string	true		"the token"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /get_exp [post]
func (w *Workers) GetExp() {
	info, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	if info.Id <= 0 {
		w.Data["json"] = ReturnError(40000, "师傅信息不存在")
		w.ServeJSON()
		return
	}
	expInfo, err := models.GetExp(info.Id)
	if err == nil && expInfo.Id > 0 {
		w.Data["json"] = ReturnSuccess(0, "success", expInfo, 1)
		w.ServeJSON()
	} else {
		w.Data["json"] = ReturnError(40001, "获取信息失败")
		w.ServeJSON()
	}
}

// @Title 经验认证提交
// @Description 师傅中心-经验认证提交接口
//// @Param	token		header 	string	true		"the token"
//// @Param	exp			query 	int 	true		"the exp"
//// @Param	wechat		query 	string 	true		"the wechat"
//// @Param	phone		query 	string 	true		"the phone"
//// @Param	image		query 	string 	true		"the image"
//// @Success 200 {string} auth success
//// @Failure 403 user not exist
// @router /save_exp [post]
func (w *Workers) SaveExp() {
	exp, _ := w.GetInt("exp")
	wechat := w.GetString("wechat")
	phone := w.GetString("phone")
	image := w.GetString("image")

	if exp <= 0 {
		w.Data["json"] = ReturnError(40000, "施工经验不能为空")
		w.ServeJSON()
		return
	}
	if wechat == "" {
		w.Data["json"] = ReturnError(40000, "微信号不能为空")
		w.ServeJSON()
		return
	}
	if phone == "" {
		w.Data["json"] = ReturnError(40000, "手机号码不能为空")
		w.ServeJSON()
		return
	}
	if image == "" {
		w.Data["json"] = ReturnError(40000, "截图证明不能为空")
		w.ServeJSON()
		return
	}
	//获取当前用户师傅id
	workerInfo, _ := models.GetMasterWorkerInfo(w.CurrentLoginUser.Id)
	//查询经验认证是否已提交
	expInfo, _ := models.GetExp(workerInfo.Id)
	if expInfo.Id > 0 {
		w.Data["json"] = ReturnError(40003, "经验认证已提交，勿重复提交")
		w.ServeJSON()
		return
	}
	boolVal, err := models.AddExp(workerInfo.Id, exp, wechat, phone, image)
	if err == nil && boolVal {
		w.Data["json"] = ReturnSuccess(0, "success", "", 1)
		w.ServeJSON()
	} else {
		w.Data["json"] = ReturnError(40001, "提交失败，请稍后再试")
		w.ServeJSON()
	}
}

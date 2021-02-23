package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

type Orders struct {
	Base
}

var constructionData = []string{"主料+辅料+施工", "仅施工", "辅料+施工"}

// @Title 订单列表
// @Description 用户订单列表
// @Param	token		header 	string	true		"the token"
// @Param	status		query 	int 	true		"the orders status"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_orders [post]
func (this *Orders) Index() {
	// 接收订单状态参数
	status, _ := this.GetInt("status")
	s := 0
	if status != 0 {
		s = status
	}
	num, data, err := models.GetOrdersAll(s)
	if err == nil && num > 0 {
		this.Data["json"] = ReturnSuccess(0, "success", data, num)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 订单详情页
// @Description 订单详情页接口数据
// @Param	token		header 	string	true		"the token"
// @Param	order_id		query 	int 	true		"the orders id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_order_details [get]
func (this *Orders) OrderDetails() {
	ValueData := make(map[string]interface{})
	// 接收订单状态参数
	orderId, _ := this.GetInt("order_id")

	data, err := models.GetOrderInfo(orderId)
	//接口数据 截止时间| 师傅列表 | 参与人数 | 订单信息 | 基本信息 | 地址信息
	tempTime := 60 * 60 * 24 * 3
	ValueData["last_time"] = (strToUnixTime(data.CreateAt) + int64(tempTime)) - time.Now().Unix()
	if err == nil {
		//获取参与报价的师傅列表最多5人
		num, list, listErr := models.GetOrderWorkerList(data.Id)
		if listErr == nil && num > 0 {
			ValueData["master_worker_list"] = list
		} else {
			ValueData["master_worker_list"] = nil
			ValueData["master_worker_list_num"] = num
		}
		//获取订单详情信息
		ValueData["info"] = data
		this.Data["json"] = ReturnSuccess(0, "success", ValueData, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 下单页数据接口
// @Description 下单页初始化数据接口
// @Success 200 {string} auth success
// @Failure 403 typeId not exist
// @router /order_page [post]
func (this *Orders) OrderPages() {
	var data = make(map[string]interface{})
	//typeId,_ := h.GetInt("type_id")

	// 获取用户信息
	data["user_name"] = this.CurrentLoginUser.Nickname
	data["phone"] = this.CurrentLoginUser.Phone
	// 服务类型，
	var typesData []string
	_, typesAll, err := models.GetAllServiceType()
	if err == nil {
		for _, val := range typesAll {
			typesData = append(typesData, val.TypeName)
		}
		data["types"] = typesData
	}
	// 施工类型 Construction type
	//var constructionData []string
	//var constructionData = []string{"主料+辅料+施工", "仅施工", "辅料+施工"}
	//constructionData[0] = "主料+辅料+施工"
	//constructionData[1] = "仅施工"
	//constructionData[2] = "辅料+施工"
	data["construction"] = constructionData

	this.Data["json"] = ReturnSuccess(40000, "success", data, 1)
	this.ServeJSON()
}

// @Title 提交订单
// @Description 普通用户提交订单
// @Param	token		header 	string	true		"the token"
// @Param	address		query 	string	true		"the address"
// @Param	construction_time		query 	string	true	"the construction_time"
// @Param	types		query 	int	true		"the types"
// @Param	is_materials		query 	int	true		"the is_materials"
// @Param	area		query 	float64	true		"the area"
// @Param	is_tear_of_old_wallpaper		query 	int	true		"the is_tear_of_old_wallpaper"
// @Param	basement_membrane		query 	int	true		"the basement_membrane"
// @Param	desc		query 	string	true		"the desc"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /confirm_order [post]
func (this *Orders) SaveOrder() {
	// 获取参数
	Address := this.GetString("address")                       // 施工地址
	constructionTime := this.GetString("construction_time")    // 施工时间
	types, typesErr := this.GetInt("types")                    // 服务类型
	materials, _ := this.GetInt("is_materials")                // 是否提供物料
	area, areaErr := this.GetFloat("area")                     // 面积
	oldWallpaper, _ := this.GetInt("is_tear_of_old_wallpaper") // 是否清除旧纸
	basementMembrane, _ := this.GetInt("basement_membrane")    // 是否刷基膜
	moreDescription := this.GetString("desc")                  // 需求描述
	images := this.GetString("images")                         // 图片

	// 用户参数
	userId := this.CurrentLoginUser.Id // 当前用户id
	//phone := "15938755991"
	orderType, orderTypesErr := this.GetInt("order_type") // 订单类型
	if orderTypesErr == nil && orderType == 2 {
		orderType = 2
	}

	// 订单类型为直接选择师傅
	workerId, workerErr := this.GetInt("worker_id") // 师傅id
	if workerErr != nil {
		workerId = 0
	}

	// 参数效验
	if Address == "" {
		this.Data["json"] = ReturnError(40001, "地址不能为空")
		this.ServeJSON()
		this.StopRun()
	}

	if constructionTime == "" {
		this.Data["json"] = ReturnError(40001, "请选择施工时间")
		this.ServeJSON()
		this.StopRun()
	}

	if typesErr != nil {
		this.Data["json"] = ReturnError(40001, "服务类型参数错误")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(types)

	if areaErr != nil || area < 1 {
		this.Data["json"] = ReturnError(40001, "请选择有效的施工面积")
		this.ServeJSON()
		this.StopRun()
	}
	fmt.Println(area)
	var orderinfo models.EOrders
	//orderinfo.OrderSn = CreateRandOrderOn() // 生成订单号
	orderinfo.MId = int(userId)
	orderinfo.WorkerId = workerId
	orderinfo.Area = area
	orderinfo.ConstructionTime = int(strToUnixTime(constructionTime))
	orderinfo.BasementMembrane = basementMembrane
	orderinfo.IsTearOfOldWallpaper = oldWallpaper
	orderinfo.IsMateriel = materials
	orderinfo.MoreDescription = moreDescription
	orderinfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	orderinfo.Status = 1
	orderinfo.OrderType = orderType
	orderinfo.Images = images
	orderinfo.Address = Address
	fmt.Println(orderinfo)
	//panic("orderinfo")
	o := orm.NewOrm()
	id, errors := o.Insert(&orderinfo)
	fmt.Println(errors)
	if errors == nil {
		this.Data["json"] = ReturnSuccess(0, "success", id, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40001, "订单提交失败")
	this.ServeJSON()
}

// @Title 取消订单
// @Description 取消订单接口
// @Param	token		header 	string	true		"the token"
// @Param	order_id		query 	int 	true		"the orders id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /order_cancel [get]
func (this *Orders) OrderCancel() {
	// 接收订单状态参数
	orderId, _ := this.GetInt("order_id")

	data, err := models.GetOrderInfo(orderId)
	if err == nil {
		boolVal, errs := models.OrderCancel(data.Id, -1)
		if boolVal {
			this.Data["json"] = ReturnSuccess(0, "success", data, 1)
			this.ServeJSON()
			return
		}
		logs.Error("取消订单操作失败:" + fmt.Sprint("%s", errs))
		this.Data["json"] = ReturnError(40000, "取消订单操作失败")
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 订单大厅
// @Description 师傅订单大厅数据接口
// @Param	token		header 	string	true		"the token"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_master_orders_list [post]
func (this *Orders) MasterOrderList() {

	num, data, err := models.GetOrdersAll(1)
	if err == nil && num > 0 {
		returnVal := make([]map[string]interface{}, len(data))
		for k := 0; k < len(returnVal); k++ {
			// 一定要加下面的nil判断  否则会报错 Handler crashed with error assignment to entry in nil map  map未赋值
			if returnVal[k] == nil {
				returnVal[k] = map[string]interface{}{}
			}
			returnVal[k]["id"] = data[k].Id
			returnVal[k]["order_sn"] = "00000"
			serviceTypeObj, _ := models.GetServiceType(int64(data[k].ServiceId))
			returnVal[k]["service_type_str"] = serviceTypeObj.TypeName
			returnVal[k]["service_time"] = data[k].ConstructionTime
			returnVal[k]["area"] = data[k].Area
			returnVal[k]["service_list"] = constructionData[data[k].IsMateriel]
			returnVal[k]["create_time"] = data[k].CreateAt
			returnVal[k]["status"] = data[k].Status
		}

		this.Data["json"] = ReturnSuccess(0, "success", returnVal, num)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 订单详情
// @Description 师傅订单大厅详情数据接口
// @Param	token		header 	string	true		"the token"
// @Param	id		query 	string	true		"the order_id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /get_master_orders_info [get]
func (this *Orders) GetMasterOrdersInfo() {
	oId, _ := this.GetInt("order_id")
	if oId == 0 {
		this.Data["json"] = ReturnError(40001, "参数错误")
		this.ServeJSON()
		return
	}
	data, err := models.GetOrderInfo(oId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", data, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40002, "订单信息不存在")
	this.ServeJSON()
}

// @Title 参与报价
// @Description 师傅参与报价确认
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /participate_offer [post]
func (this *Orders) ParticipateOffer() {

}

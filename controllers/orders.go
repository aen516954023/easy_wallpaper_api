package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
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
	num, data, err := models.GetOrdersAll(this.CurrentLoginUser.Id, s, 1)
	if err == nil && num > 0 {
		returnValue := make([]map[string]interface{}, num)
		for k, v := range data {
			if returnValue[k] == nil {
				returnValue[k] = map[string]interface{}{}
			}
			returnValue[k]["Id"] = v.Id
			returnValue[k]["Status"] = v.Status
			returnValue[k]["CreateAt"] = v.CreateAt
			returnValue[k]["ServiceId"] = v.ServiceId
			eType, _ := models.GetServiceType(int64(v.ServiceId))
			returnValue[k]["ServiceName"] = eType.TypeName
			returnValue[k]["OrderSn"] = v.OrderSn
			//returnValue[k]["WorkerId"] = v.WorkerId // Todo 订单列表师傅数据展示处理
			stepInfo, _ := models.GetOrderOfStepInfo(v.Id)
			masterInfo, _ := models.GetMasterWorkerInfId(stepInfo.WId)
			returnValue[k]["WorkerId"] = stepInfo.WId
			returnValue[k]["WorkerName"] = masterInfo.Username
			returnValue[k]["AvatarImg"] = masterInfo.Image
			if v.Status <= 1 {
				returnValue[k]["EndTime"] = GetEndTime(v.CreateAt, 3) //这个3 指3小时 取数据库参数配置时间
			}
		}
		this.Data["json"] = ReturnSuccess(0, "success", returnValue, num)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40000, "没有查询到关信息")
	this.ServeJSON()
}

// @Title 订单详情页
// @Description 用户订单详情页接口数据
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
			ValueData["master_worker_list_num"] = num
		} else {
			ValueData["master_worker_list"] = nil
			ValueData["master_worker_list_num"] = 0
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
	//data["user_name"] = this.CurrentLoginUser.Nickname
	//data["phone"] = this.CurrentLoginUser.Phone
	// 服务类型，
	var typesData []string
	_, typesAll, err := models.GetAllServiceType()
	if err == nil {
		for _, val := range typesAll {
			typesData = append(typesData, val.TypeName)
		}
		data["types"] = typesData
	}
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
	Address, _ := this.GetInt("address")                    // 施工地址
	city := this.GetString("city")                          // 施工城市
	constructionTime := this.GetString("construction_time") // 施工时间
	types, typesErr := this.GetInt("types")                 // 服务类型
	constructionType, _ := this.GetInt("construction_type") // 施工类型

	materials, _ := this.GetInt("is_materials")                // 是否提供物料
	area, areaErr := this.GetFloat("area")                     // 面积
	oldWallpaper, _ := this.GetInt("is_tear_of_old_wallpaper") // 是否清除旧纸
	basementMembrane, _ := this.GetInt("basement_membrane")    // 是否刷基膜
	moreDescription := this.GetString("desc")                  // 需求描述
	images := this.GetString("images")                         // 图片

	// 用户参数
	userId := this.CurrentLoginUser.Id                    // 当前用户id
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
	if Address == 0 {
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
	//fmt.Println(types)

	if areaErr != nil || area < 1 {
		this.Data["json"] = ReturnError(40001, "请选择有效的施工面积")
		this.ServeJSON()
		this.StopRun()
	}

	//  处理施工时间格式
	constructionTime = strings.Replace(constructionTime, "/", "-", -1)

	var orderInfo models.EOrders
	orderInfo.OrderSn = CreateRandOrderOn() // 生成订单号
	orderInfo.MId = int(userId)
	orderInfo.WorkerId = workerId
	orderInfo.Area = area
	orderInfo.ConstructionTime = int(strToUnixTime(constructionTime))
	orderInfo.BasementMembrane = basementMembrane
	orderInfo.IsTearOfOldWallpaper = oldWallpaper
	orderInfo.IsMateriel = materials
	orderInfo.MoreDescription = moreDescription
	orderInfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	orderInfo.Status = 1
	orderInfo.OrderType = orderType
	orderInfo.Images = images
	orderInfo.Address = Address
	orderInfo.City = city
	orderInfo.ConstructionType = constructionType
	orderInfo.ServiceId = types
	fmt.Println(orderInfo)
	//panic("orderinfo")
	o := orm.NewOrm()
	id, errors := o.Insert(&orderInfo)
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

	num, data, err := models.GetOrdersAll(this.CurrentLoginUser.Id, 1, 2)
	if err == nil && num > 0 {
		returnVal := make([]map[string]interface{}, len(data))
		for k := 0; k < len(returnVal); k++ {
			// 一定要加下面的nil判断  否则会报错 Handler crashed with error assignment to entry in nil map  map未赋值
			if returnVal[k] == nil {
				returnVal[k] = map[string]interface{}{}
			}
			returnVal[k]["id"] = data[k].Id
			returnVal[k]["order_sn"] = data[k].OrderSn
			serviceTypeObj, _ := models.GetServiceType(int64(data[k].ServiceId))
			returnVal[k]["service_type_str"] = serviceTypeObj.TypeName
			returnVal[k]["service_time"] = UnixTimeToSTr(int64(data[k].ConstructionTime))
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
	fmt.Println("data", data)
	if err == nil {
		returnValue := make(map[string]interface{})
		// 订单详情数据
		serviceType, _ := models.GetServiceType(int64(data.ServiceId))
		returnValue["id"] = data.Id
		returnValue["sn"] = data.OrderSn
		returnValue["service_type"] = serviceType.TypeName
		returnValue["construction_type"] = constructionData[data.ConstructionType]
		returnValue["address"] = data.Address
		returnValue["area"] = data.Area
		returnValue["create_at"] = data.CreateAt
		if data.ConstructionTime == 0 {
			returnValue["construction_time"] = "--"
		} else {
			returnValue["construction_time"] = UnixTimeToSTr(int64(data.ConstructionTime))
		}
		returnValue["basement_membrane"] = data.BasementMembrane
		returnValue["is_masteriel"] = data.IsMateriel
		returnValue["is_tear_of_old_wallpaper"] = data.IsTearOfOldWallpaper
		returnValue["desc"] = data.MoreDescription
		returnValue["images"] = data.Images

		// 查询师傅是否已经参与
		if models.GetOrderTaskingUid(oId, int(this.CurrentLoginUser.Id)) {
			returnValue["isJoin"] = 1
			// 获取电话号码
			returnValue["phone"] = this.CurrentLoginUser.Phone
		} else {
			returnValue["isJoin"] = 0
			returnValue["phone"] = ""
		}
		//通过uid查询师傅id
		masterInfo, mErr := models.GetMasterWorkerInfo(this.CurrentLoginUser.Id)
		if mErr == nil {
			returnValue["w_id"] = masterInfo.Id
		} else {
			returnValue["w_id"] = 0
		}
		this.Data["json"] = ReturnSuccess(0, "success", returnValue, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40002, "订单信息不存在")
	this.ServeJSON()
}

// @Title 参与报价
// @Description 师傅参与报价确认
// @Param	order_id		query 	int	true		"the order id"
// @Param	master_worker_id		query 	int	true		"the master worker id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /participate_offer [post]
func (this *Orders) ParticipateOffer() {
	orderId, _ := this.GetInt("o_id")
	wId, _ := this.GetInt("w_id")

	if orderId == 0 {
		this.Data["json"] = ReturnError(40001, "缺少订单id参数")
		this.ServeJSON()
		return
	}
	if wId == 0 {
		this.Data["json"] = ReturnError(40001, "缺少师傅id参数")
		this.ServeJSON()
		return
	}
	//判断师傅是否参与过
	if models.GetOrderTasking(orderId, wId) {
		this.Data["json"] = ReturnError(40002, "你已参与过了")
		this.ServeJSON()
		return
	}

	//一个订单最多参与人次
	num, countErr := models.GetTaskCount(orderId)
	if countErr == nil {
		if num >= 5 {
			this.Data["json"] = ReturnError(40002, "订单参与人数已满")
			this.ServeJSON()
			return
		}
	}

	boolVal, err := models.InsertOrderTaking(orderId, int(this.CurrentLoginUser.Id), wId)
	if boolVal {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
		return
	}
	logs.Error("师傅参与错误:" + fmt.Sprintf("%v", err))
	this.Data["json"] = ReturnError(40002, "参与失败，请稍后再试")
	this.ServeJSON()
}

// @Title 订单管理
// @Description 订单管理页面数据--师傅端
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /order_manage_master [post]
func (this *Orders) OrderManageMaster() {
	orderId, _ := this.GetInt("order_id")
	if orderId == 0 {
		this.Data["json"] = ReturnError(40001, "参数错误，缺少订单id")
		this.ServeJSON()
		return
	}
	returnVal := make(map[string]interface{})
	masterInfo, _ := models.GetMasterWorkerInfo(this.CurrentLoginUser.Id)
	data, err := models.GetOrderOfStepOne(orderId, masterInfo.Id, true)
	if err == nil && data.Id > 0 {
		returnVal["WId"] = data.WId
		wInfo, _ := models.GetMasterWorkerInfId(data.WId)
		fmt.Println("wInfo", wInfo)
		returnVal["w_username"] = wInfo.Username
		returnVal["w_avatar"] = wInfo.Image
		mInfo, _ := models.GetMemberInfoId(wInfo.MId)
		returnVal["phone"] = mInfo.Phone
		oInfo, _ := models.GetOrderInfo(data.OId)
		returnVal["o_status"] = oInfo.Status
		returnVal["Area"] = data.Area
		returnVal["ConstructionType"] = data.ConstructionType
		returnVal["ConstructionTypeStr"] = constructionData[data.ConstructionType]
		returnVal["CreateAt"] = data.CreateAt
		//DepositId: 0
		returnVal["DepositPrice"] = data.DepositPrice
		returnVal["DepositStatus"] = data.DepositStatus
		returnVal["DiscountedPrice"] = data.DiscountedPrice
		returnVal["HomeTime"] = data.HomeTime
		returnVal["Id"] = data.Id
		returnVal["Info"] = data.Info
		returnVal["MId"] = data.MId
		returnVal["OId"] = data.OId
		//PayId: 0
		returnVal["PayStatus"] = data.PayStatus
		returnVal["Price"] = data.Price
		sInfo, _ := models.GetServiceType(int64(data.ServiceType))
		returnVal["ServiceType"] = data.ServiceType
		returnVal["ServiceTypeStr"] = sInfo.TypeName
		returnVal["Step1"] = data.Step1
		returnVal["Step1Time"] = UnixTimeToSTr(int64(data.Step1Time))
		returnVal["Step2"] = data.Step2
		returnVal["Step2Time"] = UnixTimeToSTr(int64(data.Step2Time))
		returnVal["Step3"] = data.Step3
		returnVal["Step3Time"] = UnixTimeToSTr(int64(data.Step3Time))
		returnVal["Step4"] = data.Step4
		returnVal["Step4Time"] = UnixTimeToSTr(int64(data.Step4Time))
		returnVal["Step5"] = data.Step5
		returnVal["Step5Time"] = UnixTimeToSTr(int64(data.Step5Time))
		returnVal["Step6"] = data.Step6
		returnVal["Step6Time"] = UnixTimeToSTr(int64(data.Step6Time))
		returnVal["Step7"] = data.Step7
		returnVal["Step7Time"] = UnixTimeToSTr(int64(data.Step7Time))
		returnVal["TotalPrice"] = data.TotalPrice
		returnVal["Unit"] = data.Unit

		_, list, _ := models.GetAllServiceType()
		returnVal["service_type_list"] = list
		returnVal["construction_type_list"] = constructionData
		this.Data["json"] = ReturnSuccess(0, "success", returnVal, 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40001, "暂无记录")
		this.ServeJSON()
	}
}

// @Title 订单管理
// @Description 订单管理页面数据--客户端
// @Param	order_id		query 	int	true		"the order id"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /order_manage_user [post]
func (this *Orders) OrderManageUser() {
	orderId, _ := this.GetInt("order_id")
	if orderId == 0 {
		this.Data["json"] = ReturnError(40001, "参数错误，缺少订单id")
		this.ServeJSON()
		return
	}
	returnVal := make(map[string]interface{})
	data, err := models.GetOrderOfStepOne(orderId, int(this.CurrentLoginUser.Id), false)
	if err == nil && data.Id > 0 {
		returnVal["WId"] = data.WId
		wInfo, _ := models.GetMasterWorkerInfId(data.WId)
		fmt.Println("wInfo", wInfo)
		returnVal["w_username"] = wInfo.Username
		returnVal["w_avatar"] = wInfo.Image
		mInfo, _ := models.GetMemberInfoId(wInfo.MId)
		returnVal["phone"] = mInfo.Phone
		oInfo, _ := models.GetOrderInfo(data.OId)
		returnVal["o_status"] = oInfo.Status
		returnVal["Area"] = data.Area
		returnVal["ConstructionType"] = constructionData[data.ConstructionType]
		returnVal["CreateAt"] = data.CreateAt
		//DepositId: 0
		returnVal["DepositPrice"] = data.DepositPrice
		returnVal["DepositStatus"] = data.DepositStatus
		returnVal["DiscountedPrice"] = data.DiscountedPrice
		returnVal["HomeTime"] = UnixTimeToSTr(int64(data.HomeTime))
		returnVal["Id"] = data.Id
		returnVal["Info"] = data.Info
		returnVal["MId"] = data.MId
		returnVal["OId"] = data.OId
		//PayId: 0
		returnVal["PayStatus"] = data.PayStatus
		returnVal["Price"] = data.Price
		sInfo, _ := models.GetServiceType(int64(data.ServiceType))
		returnVal["ServiceType"] = sInfo.TypeName
		returnVal["Step1"] = data.Step1
		returnVal["Step1Time"] = UnixTimeToSTr(int64(data.Step1Time))
		returnVal["Step2"] = data.Step2
		returnVal["Step2Time"] = UnixTimeToSTr(int64(data.Step2Time))
		returnVal["Step3"] = data.Step3
		returnVal["Step3Time"] = UnixTimeToSTr(int64(data.Step3Time))
		returnVal["Step4"] = data.Step4
		returnVal["Step4Time"] = UnixTimeToSTr(int64(data.Step4Time))
		returnVal["Step5"] = data.Step5
		returnVal["Step5Time"] = UnixTimeToSTr(int64(data.Step5Time))
		returnVal["Step6"] = data.Step6
		returnVal["Step6Time"] = UnixTimeToSTr(int64(data.Step6Time))
		returnVal["Step7"] = data.Step7
		returnVal["Step7Time"] = UnixTimeToSTr(int64(data.Step7Time))
		returnVal["TotalPrice"] = data.TotalPrice
		returnVal["Unit"] = data.Unit

		this.Data["json"] = ReturnSuccess(0, "success", returnVal, 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40001, "暂无记录")
		this.ServeJSON()
	}
}

package controllers

import (
	"easy_wallpaper_api/models"
	"strconv"
	"time"
)

type User struct {
	Base
}

// @Title 个人中心
// @Description 个人中心页数据接口
// @Param	token header string	true "token值"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /my/index [post]
func (this *User) Index() {
	//获取当前时间戳
	currentTime := time.Now().Unix()
	infoMap := make(map[string]interface{})
	//获取登录用户的信息
	uid := this.CurrentLoginUser.Id
	//获取总订单数量
	order_all_num, _ := models.GetOrderAllNum(uid)
	//获取生效订单数量
	order_effective_num, total_power, _ := models.GetOrderEffectiveNum(uid, currentTime)
	//获取btc兑换美金的比例
	btcU, _ := strconv.ParseFloat(Config("btc_u"), 64)
	infoMap["phone"] = this.CurrentLoginUser.Name
	infoMap["last_login_time"] = time.Unix(this.CurrentLoginUser.LastLoginTime, 0).Format("2006-01-02 03:04:05 ")
	infoMap["user_img"] = "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1851283359,3457678391&fm=26&gp=0.jpg"
	infoMap["order_num"] = order_all_num                           // 订单数量
	infoMap["account_fee"] = 0                                     // 待交电费
	infoMap["card_num"] = 0                                        // 卡包
	infoMap["total_power"] = total_power                           // 总算力
	infoMap["take_effecting"] = order_effective_num                // 生效中
	infoMap["total_historical_output"] = this.CurrentLoginUser.Btc // 历史总产出
	infoMap["about_us_dollars"] = this.CurrentLoginUser.Btc * btcU //约等于多少美金
	//获取当前用户前七天的收益
	countryCapitalMap, boolVal, err := models.GetcountryCapital(uid)
	if err == nil && boolVal {
		infoMap["week_data"] = countryCapitalMap //一周产出图表数据
	} else {
		infoMap["week_data"] = 0 //一周产出图表数据
	}

	this.Data["json"] = ReturnSuccess(0, "success", infoMap, 1)
	this.ServeJSON()
}

// @Title 我的订单
// @Description 个人中心-我的订单
// @Param	token header string	true "token值"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /my/orders [post]
func (this *User) OrderList() {
	nums, infos, err := models.GetAllOrders(this.CurrentLoginUser.Id)
	if err == nil {
		infoMaps := make([]map[string]interface{}, len(infos))
		for k := 0; k < len(infoMaps); k++ {
			// 一定要加下面的nil判断  否则会报错 Handler crashed with error assignment to entry in nil map  map未赋值
			if infoMaps[k] == nil {
				infoMaps[k] = map[string]interface{}{}
			}
			infoMaps[k]["id"] = infos[k].Id
			infoMaps[k]["name"] = infos[k].GoodsSkuName
			infoMaps[k]["status"] = infos[k].OrderStatus
			infoMaps[k]["time"] = infos[k].CreateAt
			infoMaps[k]["total_price"] = infos[k].TotalPrice
			infoMaps[k]["Transit_price"] = infos[k].TransitPrice
			infoMaps[k]["nums"] = infos[k].TotalCount
		}
		this.Data["json"] = ReturnSuccess(0, "success", infoMaps, nums)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40001, "未查询到相关订单信息")
		this.ServeJSON()
	}

}

// @Title 个人中心
// @Description 个人中心页数据接口
// @Param	token header string	true "token值"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /my/earnlist [post]
func (this *User) Earnlist() {
	//获取登录用户的信息
	uid := this.CurrentLoginUser.Id
	//获取我的收益列表
	num, earnList, err := models.GetEarnList(uid)
	if err == nil {
		datalist := make([]map[string]interface{}, num)
		for key, value := range earnList {
			datalist[key] = map[string]interface{}{}
			datalist[key]["earn_num"] = value.Btc
			datalist[key]["des"] = value.Des
			datalist[key]["add_time"] = time.Unix(value.AddTime, 0).Format("2006-01-02 03:04:05 ")
		}
		this.Data["json"] = ReturnSuccess(0, "success", datalist, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(1, "error")
		this.ServeJSON()
	}
}

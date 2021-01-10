package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type Home struct {
	beego.Controller
}

// @Title 首页数据接口
// @Description 首页数据接口
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /index [get]
func (h *Home) Index() {
	var data = make(map[string]interface{})
	// banner
	_, banner, err := models.GetBannerAll()
	if err == nil {
		data["banner"] = banner
	} else {
		data["banner"] = nil
		//fmt.Println(err)
	}
	// msg
	num, msg, errs := models.GetMsgBefore(3)
	if errs == nil && num > 0 {
		data["msg"] = msg
	} else {
		data["msg"] = nil
		//fmt.Println(errs)
	}
	// 推荐列表
	numg, goods, errg := models.GetRecommendList(4)
	if errg == nil && numg > 0 {
		// 处理数据转换
		goodList := make([]map[string]interface{}, len(goods))
		for k := 0; k < len(goods); k++ {
			// 一定要加下面的nil判断  否则会报错 Handler crashed with error assignment to entry in nil map  map未赋值
			if goodList[k] == nil {
				goodList[k] = map[string]interface{}{}
			}
			goodList[k]["id"] = goods[k].Id
			goodList[k]["title"] = goods[k].Goods.Name
			goodList[k]["effect_time"] = goods[k].BeginTime
			goodList[k]["price"] = goods[k].Price
			goodList[k]["cycle"] = goods[k].GoodsCycle.Day
			goodList[k]["power"] = goods[k].Goods.Power
			goodList[k]["stock"] = goods[k].Stock
			goodList[k]["sales"] = goods[k].Sales
			progress, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(goods[k].Sales)/float64(goods[k].Stock)*100), 64)
			goodList[k]["progress"] = progress

			if goods[k].Sales >= goods[k].Stock {
				goodList[k]["sold"] = 1
			} else {
				goodList[k]["sold"] = 0
			}
			goodList[k]["static_income"] = goods[k].StaticIncome
			goodList[k]["power_cost"] = goods[k].Goods.Hashrate
			goodList[k]["static_output"] = goods[k].Goods.Hashrate
			goodList[k]["electricity_fees"] = goods[k].Goods.Hashrate
			goodList[k]["type"] = goods[k].GoodsMode.Id
			goodList[k]["type_name"] = goods[k].GoodsMode.Name
		}
		data["goods"] = goodList
	} else {
		data["goods"] = nil
		fmt.Println(errg)
	}

	h.Data["json"] = ReturnSuccess(0, "success", data, 0)
	h.ServeJSON()
}

// @Title 币种数据接口
// @Description 币种数据接口
// @Success 200 {string} auth success
// @Failure 403 typeId not exist
// @router /type_push [get]
func (h *Home) TypeDataPush() {
	num, data, err := models.GetAllCurrency()
	if err == nil && num > 0 {
		flag, _ := h.GetInt("flag")
		if flag == 1 {
			dataMap := make(map[int]string, len(data))
			for i := 0; i < len(data); i++ {
				//if dataMap[i] == nil {
				//	dataMap[i] = map[int]string
				//}
				dataMap[data[i].Id] = data[i].Name
			}
			h.Data["json"] = ReturnSuccess(0, "success", dataMap, 0)
		} else {
			h.Data["json"] = ReturnSuccess(0, "success", data, 0)
		}
		h.ServeJSON()
	} else {
		h.Data["json"] = ReturnError(40000, "error")
		h.ServeJSON()
	}
}

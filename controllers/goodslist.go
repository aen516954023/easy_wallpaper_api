package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type GoodsList struct {
	beego.Controller
}

// @Title 列表页
// @Description 列表页接口
// @Success 200 {string} auth success
// @Failure 403
// @router /Index [post]
func (this *GoodsList) Index() {
	//获取所有币种列表，
	num, type_data, err := models.GetAllCurrency()
	if err != nil {
		//查询出错
		this.Data["json"] = ReturnError(0, "获取币种列表出错~")
		this.ServeJSON()
		this.StopRun()
	}
	//获取所有矿机列表，
	num, kj_data, err := models.GetAllkj()
	if err != nil {
		//查询出错
		this.Data["json"] = ReturnError(0, "获取矿机列表出错~")
		this.ServeJSON()
		this.StopRun()
	}
	//获取所有周期列表，
	num, zq_data, err := models.GetAllzq()
	if err != nil {
		//查询出错
		this.Data["json"] = ReturnError(0, "获取周期列表出错~")
		this.ServeJSON()
		this.StopRun()
	}
	//获取所有套餐模式列表，
	num, tc_data, err := models.GetAlltc()
	if err != nil {
		//查询出错
		this.Data["json"] = ReturnError(0, "获取套餐模式列表出错~")
		this.ServeJSON()
		this.StopRun()
	}

	//获取默认为币种为BTC的矿机列表
	//获取第一个币种类型的id
	type_value, err := models.Get_first_type_id()
	if err != nil {
		//查询出错
		this.Data["json"] = ReturnError(0, "获取币种出错~")
		this.ServeJSON()
		this.StopRun()
	}
	//获取默认为币种为BTC的矿机列表
	goodskeo := models.Get_goods_by_type(type_value.Id)
	//返回币种列表
	msg := map[string]interface{}{
		"type_data": type_data,
		"kj_data":   kj_data,
		"tc_data":   tc_data,
		"zq_data":   zq_data,
		"goodskeo":  goodskeo,
	}

	this.Data["json"] = ReturnSuccess(0, "success", msg, num)
	this.ServeJSON()
	this.StopRun()

}

// @Title 根据条件查询列表信息
// @Description 根据条件查询列表信息
//@Param    type    query    string         false        "币种id"
//@Param    cycle    query    string         false        "周期id"
//@Param    mode    query    string         false        "模式id"
//@Param    kjid    query    string         false        "矿机id"
// @Success 200 {string} auth success
// @Failure 403
// @router /GetGoodslist [post]
func (this *GoodsList) GetGoodsList() {
	where := make(map[string]string)
	type_value := this.GetString("type")
	if type_value != "" {
		where["goods_type_id"] = type_value
	}
	cycle_value := this.GetString("cycle")
	if cycle_value != "" {
		where["goods_cycle_id"] = cycle_value
	}
	mode_value := this.GetString("mode")
	if mode_value != "" {
		where["goods_mode_id"] = mode_value
	}
	kj_value := this.GetString("kjid")
	if mode_value != "" {
		where["goods_id"] = kj_value
	}
	goods := models.GetGoodsListByTj(where)
	if len(goods) > 0 {
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
		this.Data["json"] = ReturnSuccess(0, "success", goodList, int64(len(goodList)))
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(40001, "没有查询到相关数据")
		this.ServeJSON()
	}
}

// @Title 详情页
// @Description goods详情页
//@Param    goodsskuid    query    int         true        "goodssku的id"
// @Success 200 {string} auth success
// @Failure 403
// @router /GetGoodsdetail [post]
func (this *GoodsList) GetGoodsDetail() {
	goodsskuid, err := this.GetInt("goodsskuid")
	if err != nil {
		fmt.Printf("222", err)
		//查询出错
		this.Data["json"] = ReturnError(0, "请输入产品id~")
		this.ServeJSON()
		this.StopRun()
	}
	//获取默认为币种为BTC的矿机列表
	goodskeo := models.Get_goods_by_goodsskuid(goodsskuid)
	this.Data["json"] = ReturnSuccess(0, "success", goodskeo, 1)
	this.ServeJSON()
	this.StopRun()
}

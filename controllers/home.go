package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
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
	// 分类
	num, typeData, typeErr := models.GetAllServiceType()
	if typeErr == nil && num > 0 {
		data["types"] = typeData
	} else {
		data["types"] = ""
	}

	// 师傅推荐列表
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
			goodList[k]["name"] = "王师傅"
			goodList[k]["effect_time"] = 333
			goodList[k]["price"] = 333
			//progress, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(goods[k].Sales)/float64(goods[k].Stock)*100), 64)
			//goodList[k]["progress"] = progress
			//
			//if goods[k].Sales >= goods[k].Stock {
			//	goodList[k]["sold"] = 1
			//} else {
			//	goodList[k]["sold"] = 0
			//}
			//goodList[k]["static_income"] = goods[k].StaticIncome
			//goodList[k]["power_cost"] = goods[k].Goods.Hashrate
			//goodList[k]["static_output"] = goods[k].Goods.Hashrate
			//goodList[k]["electricity_fees"] = goods[k].Goods.Hashrate
			//goodList[k]["type"] = goods[k].GoodsMode.Id
			//goodList[k]["type_name"] = goods[k].GoodsMode.Name
		}
		data["goods"] = goodList
	} else {
		data["goods"] = nil
		fmt.Println(errg)
	}

	h.Data["json"] = ReturnSuccess(0, "success", data, 0)
	h.ServeJSON()
}

// @Title 图片上传
// @Description 图片上传接口
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /uploads [post]
func (this *Home) Uploads() {

	f, h, err := this.GetFile("image")
	if err != nil {
		logs.Error("getfile err ", err)
		this.Data["json"] = ReturnError(40001, "上传失败")
		this.ServeJSON()
		return
	}
	defer f.Close()
	fileName := time.Now().Format("2006-01-02-15-04-05-") + h.Filename
	path := fmt.Sprint("static/upload/", fileName)
	fmt.Println(fileName)
	errs := this.SaveToFile("image", path) // 保存位置在 static/upload, 没有文件夹要先创建
	if errs == nil {
		this.Data["json"] = ReturnSuccess(0, "message", beego.AppConfig.String("appurl")+"/"+path, 1)
		this.ServeJSON()
		return
	}
	logs.Error("uploads err ", errs)
	this.Data["json"] = ReturnError(40003, "上传失败")
	this.ServeJSON()
}

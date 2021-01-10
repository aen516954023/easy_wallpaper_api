package controllers

import (
	"easy_wallpaper_api/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"strconv"
	"time"
)

func InitTask() {

	tk := toolbox.NewTask("cancelOrder", "0 */1 * * * *", CancelOrder)
	tl := toolbox.NewTask("sendFee", "0 */2 * * * *", SendFee)
	toolbox.AddTask("cancelOrder", tk)
	toolbox.AddTask("sendFee", tl)
}

// 定时任务-处理订单超时 取消30分钟内未支付的订单
func CancelOrder() error {
	//获取未支付的订单数据
	num, list, err := models.GetUnpaidOrders()
	if err != nil && num <= 0 {
		return nil
	}
	// 比对下单时间是否超时
	for _, item := range list {
		currentTime := strToUnixTime(item.CreateAt)
		if (time.Now().Unix() - 30*60) > currentTime {
			// 超时 更新订单状态
			o := orm.NewOrm()
			num, err := o.QueryTable("order_info").Filter("id", item.Id).Update(orm.Params{
				"order_status": -1,
			})
			if err != nil || num == 0 {
				logs.Error("当前订单处理超时错误：" + fmt.Sprintf("%s", err))
			}
		}
		//未超时不做处理
	}
	return nil
}

//定时发放收益
func SendFee() error {
	//获取当前时间戳
	currentTime := time.Now().Unix()
	//获取符合发放条件的订单
	allorders := models.GetAllRorders(currentTime)
	//循环按照条件发放收益，并记录
	//循环按照条件发放收益，并记录
	for _, item := range allorders {
		//获取用户uid,增加用户的收益
		uid := item.User.Id
		//获取比例
		slbi, _ := strconv.ParseFloat(Config("fee_btc"), 64)
		ysy := item.User.Btc
		allsy := float64(item.Unite) * slbi
		//增加用户收益,增加记录
		models.AddUserFee(uid, allsy, ysy)
	}
	return nil
}

//定时获取非小号中前十位的币种的价格
//func getfxh()error{
//
//
//
//
//
//}

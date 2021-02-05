package controllers

func InitTask() {

	//tk := toolbox.NewTask("cancelOrder", "0 */1 * * * *", CancelOrder)
	//tl := toolbox.NewTask("sendFee", "0 */2 * * * *", SendFee)
	//toolbox.AddTask("cancelOrder", tk)
	//toolbox.AddTask("sendFee", tl)
}

// 定时任务-处理订单超时 取消30分钟内未支付的订单
func CancelOrder() error {
	////获取未支付的订单数据
	//num, list, err := models.GetUnpaidOrders()
	//if err != nil && num <= 0 {
	//	return nil
	//}
	//// 比对下单时间是否超时
	//for _, item := range list {
	//	currentTime := strToUnixTime(item.CreateAt)
	//	if (time.Now().Unix() - 30*60) > currentTime {
	//		// 超时 更新订单状态
	//		o := orm.NewOrm()
	//		num, err := o.QueryTable("order_info").Filter("id", item.Id).Update(orm.Params{
	//			"order_status": -1,
	//		})
	//		if err != nil || num == 0 {
	//			logs.Error("当前订单处理超时错误：" + fmt.Sprintf("%s", err))
	//		}
	//	}
	//	//未超时不做处理
	//}
	return nil
}

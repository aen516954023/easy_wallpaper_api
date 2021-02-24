package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/core/validation"
	"log"
)

type Test struct {
	beego.Controller
}

// @Title 测试接口
// @Description
// @Param	order_id		query 	int	true		"the order id"
// @Param	pay_type		query 	int	true		"the pay type"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /test [post]
//
func (this *Test) Index() {
	// 接收参数
	orderId, _ := this.GetInt("order_id")
	//mId := 1
	//wId, _ := this.GetInt("w_id")
	//serviceType, _ := this.GetInt("service_type")
	//constructionType, _ := this.GetInt("construction_type")
	//price,_ := this.GetFloat("price")
	//unit,_ := this.GetInt("unit")
	//info := this.GetString("info")
	//depositPrice,_ := this.GetFloat("deposit_price")
	// 参数效验 Todo
	valid := validation.Validation{}
	valid.Required(orderId, "order_id")
	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Fatal(err.Key, err.Message)
		}
		return
	}
	log.Fatal("success")
}

package controllers

import (
	"easy_wallpaper_api/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/xingliuhua/leaf"
	"math/rand"
	"time"
)

type CommonController struct {
	beego.Controller
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg interface{}, items interface{}, count int64) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg, Items: items, Count: count}
	return
}
func ReturnError(code int, msg interface{}) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg}
	return
}
func CreateCaptcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// 生成订单号
func CreateRandOrderOn() string {
	err, node := leaf.NewNode(0)
	if err != nil {
		return ""
	}
	err, id := node.NextId()
	if err != nil {
		return ""
	}
	return time.Now().Format("0102150405") + id
}

// 时间转时间戳
func strToUnixTime(str string) int64 {
	tmp := "2006-01-02 15:04:05"
	res, _ := time.ParseInLocation(tmp, str, time.Local)
	return res.Unix()
}

// 时间戳转时间
func UnixTimeToSTr(timestamp int64) string {
	objectTime := time.Unix(timestamp, 0)
	date := objectTime.Format("2006-01-02 15:04:05")
	return date
}

// 获取config参数
func Config(field string) string {
	if field == "" {
		return ""
	}
	data, _ := models.GetConfig(field)
	return data.Value
}

//json转化为字典
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}
	return m, nil
}

//获取支付订单，跳转url
func GetOrderUrl(data map[string]interface{}) map[string]interface{} {
	postdata := make(map[string]interface{})
	postdata["orders_code"] = data["orders"]              //订单号
	postdata["order_total"] = data["order_total"]         //支付总金额
	postdata["currency_code"] = data["currency_code"]     //币种，例：美金USD
	postdata["order_total_usd"] = data["order_total_usd"] //总折算美金金额
	postdata["notify_url"] = data["notify_url"]           //支付结果回调地址
	productdata := make(map[string]interface{})
	productdata["products_id"] = data["products_id"]               //产品id
	productdata["products_name"] = data["products_name"]           //产品名称
	productdata["products_price"] = data["products_price"]         //产品价格
	productdata["products_price_usd"] = data["products_price_usd"] //产品折算美金价格
	postdata["user_id"] = 5588
	postdata["keys"] = "y4utxyNQNHpGpQT0"
	postdata["domain"] = "http://42.51.10.75:7999"
	postdata["pay_code"] = "m_stripe"
	postdata["first_name"] = "btxl"
	postdata["customers_email"] = "1017093063@qq.com"
	postdata["street_address"] = "qingdao"
	postdata["city"] = "qingdao"
	postdata["state"] = "qingdao"
	postdata["country"] = "china"
	postdata["country_code"] = "457"
	postdata["postcode"] = "430035"
	postdata["customers_telephone"] = "17332123212"
	postdata["shipping"] = "0.00"
	postdata["cancel_url"] = ""
	postdata["return_url"] = ""
	postdata["ip"] = "127.0.0.1"
	productdata["products_quantity"] = 1
	productdata["products_image"] = "111"
	productdata["products_attributes"] = "1"
	productdata["products_model"] = "1"
	postdata["products"] = productdata
	sumbit_data, err := json.MarshalIndent(postdata, "", " ")
	if err != nil {
		fmt.Println("json.Marshal error")
	}
	url := "https://paydongfang.com/query.php"
	req := httplib.Post(url)
	req.Param("data", string(sumbit_data))
	resa, _ := req.String()
	mapj, _ := JsonToMap(resa)
	resultdata := make(map[string]interface{})
	resultdata["orders_id"] = mapj["data"].(map[string]interface{})["orders_id"]
	resultdata["fetch_url"] = "https://paydongfang.com/stripe_fetch.php"
	resultdata["callback_url"] = "https://paydongfang.com/stripe_callback.php"
	resultdata["redirect_url"] = mapj["data"].(map[string]interface{})["redirect_url"].(string) + mapj["data"].(map[string]interface{})["redirect_name"].(string)
	return resultdata
}

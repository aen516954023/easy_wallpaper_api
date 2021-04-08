package controllers

import (
	"easy_wallpaper_api/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/xingliuhua/leaf"
	"math/rand"
	"strconv"
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

//订单结束时间计算 cTime 订单创建时间， sTime 后台参数设置的时间
func GetEndTime(cTime string, hour int) int64 {
	tempTime := 60 * 60 * 24 * hour // 一天的秒数
	return (strToUnixTime(cTime) + int64(tempTime)) - time.Now().Unix()
}

// float64转int64
func Float64ToInt64(s float64) int64 {
	str := strconv.FormatFloat(s, 'E', -1, 64)
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

//获取小程序access_token
func GetAccessToken() {
	appId := beego.AppConfig.String("appId")
	secret := beego.AppConfig.String("secret")
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + secret
	req := httplib.Get(url)
	var result map[string]interface{}
	req.ToJSON(&result)
	expiresIn := Float64ToInt64(result["expires_in"].(float64))
	bm.Put("access_token", result["access_token"], time.Duration(expiresIn)*time.Second)
	//fmt.Println("access_token",bm.Get("access_token"))
}

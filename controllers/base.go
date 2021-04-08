package controllers

import (
	"easy_wallpaper_api/models"
	"easy_wallpaper_api/util"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"strings"
)

type Base struct {
	beego.Controller
	CurrentLoginUser models.EMembers
}

func (this *Base) Prepare() {
	this.Auth()
}

// 权限验证
func (Base *Base) Auth() {
	token := Base.Ctx.Request.Header.Get("Token")
	if token == "" {
		//token为空跳转到授权页
		Base.Data["json"] = ReturnError(10002, "token为空，跳转到授权页")
		Base.ServeJSON()
	} else {
		userInfo, err := util.ValidateToken(token)
		if err != nil {
			//解析失败
			Base.Data["json"] = ReturnError(10002, "token验证失败，跳转到授权页重新生成token")
			Base.ServeJSON()
		} else {
			//根据解析的member查询出用户基本信息
			members, err := models.GetMemberInfo(userInfo.OpenId)
			if err == nil {
				Base.CurrentLoginUser = members
			}
		}
	}
	Base.Data["user"] = Base.CurrentLoginUser
}

func (this *Base) getClientIp() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

//模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }
//定义消息模板数组
var str = [3]string{0: "7kpZ-X3eqWrsd_kJuHUOo8Q71yXtBncP1iljrohv1x0", 1: "tom", 2: ""}

// 发送订阅消息
func (w *Base) sendSubMessage(openId string, tmpId int, data map[string]interface{}) {
	// 1 刷新ACCESS_TOKEN
	GetAccessToken()
	//2 获取 ACCESS_TOKEN
	accessToken := bm.Get("access_token")
	//3 发送订阅消息
	url := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + accessToken.(string)
	//拼接参数
	jsonData := gabs.New()
	jsonData.Set(accessToken.(string), "access_token")
	jsonData.Set(openId, "touser")
	jsonData.Set(str[tmpId], "template_id")
	jsonData.Set(data["page"], "page")
	switch tmpId {
	case 0:
		jsonData.SetP(data["order_sn"], "data.character_string5.value")
		jsonData.SetP(data["master_name"], "data.thing6.value")
		jsonData.SetP(data["date_time"], "data.time3.value")
		jsonData.SetP(data["desc"], "data.thing4.value")
	case 1:
		jsonData.SetP(data["service_name"], "data.thing1.value")
		jsonData.SetP(data["status"], "data.phrase2.value")
		jsonData.SetP(data["desc"], "data.thing4.value")
	case 2:
		jsonData.SetP(data["order_sn"], "data.character_string1.value")
		jsonData.SetP(data["status"], "data.phrase2.value")
		jsonData.SetP(data["desc"], "data.thing3.value")
	}

	//发送请求
	req := httplib.Post(url)
	req.Body(jsonData.String())
	var result map[string]interface{}
	req.ToJSON(&result)
	if result["errcode"] != 0 {
		logs.Error("订阅消息发送失败:", result["errmsg"])
	}
	fmt.Println("result", result)
}

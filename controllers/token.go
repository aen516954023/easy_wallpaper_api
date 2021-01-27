package controllers

import (
	"easy_wallpaper_api/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type Token struct {
	beego.Controller
}

type TokenData struct {
	OpenId     string
	SessionKey string
}

// @Title 小程序登陆
// @Description 小程序登陆接口
// @Param	code	query 	string	true	"小程序code码"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /login [get]
func (this *Token) Login() {
	code := this.GetString("code")
	if code == "" {
		this.Data["json"] = ReturnError(40001, "code参数不能为空")
		this.ServeJSON()
		return
	}
	// 拿code换取openid
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + beego.AppConfig.String("appId") + "&secret=" + beego.AppConfig.String("secret") + "&js_code=" + code + "&grant_type=authorization_code"

	result := httplib.Get(url)
	var res TokenData
	result.ToJSON(&res)
	fmt.Println(res)
	// 查看用户是否存在 如果存在 刷新token 如果不存在 新增用户信息

	// 生成令牌 返回前端
	token, err := util.GenerateToken(res.OpenId, 0)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", token, 1)
		this.ServeJSON()
		return
	}
	this.Data["json"] = ReturnError(40003, "登陆失败")
	this.ServeJSON()
}

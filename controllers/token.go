package controllers

import (
	"easy_wallpaper_api/models"
	"easy_wallpaper_api/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

type Token struct {
	beego.Controller
}

type TokenData struct {
	OpenId     string `json"open_id"`
	SessionKey string `json:"session_key"`
}

// @Title 小程序登陆
// @Description 小程序登陆接口
// @Param	code	query 	string	true	"小程序code码"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /login [post]
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
	//fmt.Println(res.OpenId != "")
	if res.OpenId == "" {
		this.Data["json"] = ReturnError(40001, "参数错误，获取OpenId信息失败")
		this.ServeJSON()
		return
	}
	// 查看用户是否存在 如果存在 刷新token 如果不存在 新增用户信息
	info, infoErr := models.GetMemberInfo(res.OpenId)
	fmt.Println(infoErr, info.Id)
	fmt.Println(infoErr == nil)
	if infoErr != nil && info.Id == 0 {
		//创建新用户
		insertId, err := models.AddMember(res.OpenId)
		if err == nil && insertId > 0 {
			// 生成令牌 返回前端
			token, err := util.GenerateToken(&util.User{Id: insertId, OpenId: res.OpenId}, 0)
			if err == nil {
				this.Data["json"] = ReturnSuccess(0, "success", token, 1)
				this.ServeJSON()
				return
			}
		}
	} else {
		// 用户已存在，刷新Token
		token, err := util.GenerateToken(&util.User{Id: info.Id, OpenId: info.OpenId}, 0)
		if err == nil {
			this.Data["json"] = ReturnSuccess(0, "success", token, 1)
			this.ServeJSON()
			return
		}
	}
	this.Data["json"] = ReturnError(40003, "登陆失败")
	this.ServeJSON()
}

// @Title 令牌验证
// @Description 令牌验证
// @Param	code	query 	string	true	"小程序code码"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /verify [post]
func (this *Token) Verify() {

	token := this.GetString("token")
	if token == "" {
		this.Data["json"] = ReturnSuccess(0, "success", false, 0)
		this.ServeJSON()
		return
	}
	info, err := util.ValidateToken(token)
	if err == nil {
		if info.Id > 0 {
			this.Data["json"] = ReturnSuccess(0, "success", true, 0)
			this.ServeJSON()
			return
		}
	}
	this.Data["json"] = ReturnSuccess(0, "success", false, 0)
	this.ServeJSON()
}

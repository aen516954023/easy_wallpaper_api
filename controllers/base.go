package controllers

import (
	"easy_wallpaper_api/models"
	"easy_wallpaper_api/util"
	"github.com/astaxie/beego"
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

// @APIVersion 1.0.0
// @Title 易贴壁纸 API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"easy_wallpaper_api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/home",
			beego.NSInclude(
				&controllers.Home{},
			),
		),
		beego.NSNamespace("/goodslist",
			beego.NSInclude(
				&controllers.GoodsList{},
			),
		),
		beego.NSNamespace("/orders",
			beego.NSInclude(
				&controllers.Orders{},
			),
		),
		beego.NSNamespace("/token",
			beego.NSInclude(
				&controllers.Token{},
			),
		),
		beego.NSNamespace("/workers",
			beego.NSInclude(
				&controllers.Workers{},
			),
		),
		beego.NSNamespace("/notify",
			beego.NSInclude(
				&controllers.Notify{},
			),
		),
	)
	beego.AddNamespace(ns)
}

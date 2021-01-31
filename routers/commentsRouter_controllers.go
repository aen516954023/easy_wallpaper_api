package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"],
		beego.ControllerComments{
			Method:           "EmailCreate",
			Router:           "/EmailCreate",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"],
		beego.ControllerComments{
			Method:           "EmailLogin",
			Router:           "/EmailLogin",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Auth"],
		beego.ControllerComments{
			Method:           "SendEmailVer",
			Router:           "/SendEmailVer",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"],
		beego.ControllerComments{
			Method:           "GetGoodsDetail",
			Router:           "/GetGoodsdetail",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"],
		beego.ControllerComments{
			Method:           "GetGoodsList",
			Router:           "/GetGoodslist",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:GoodsList"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/Index",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Home"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Home"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/index",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Home"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Home"],
		beego.ControllerComments{
			Method:           "TypeDataPush",
			Router:           "/type_push",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Notify"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Notify"],
		beego.ControllerComments{
			Method:           "CallbackNotify",
			Router:           "/call_back",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "SaveOrder",
			Router:           "/confirm_order",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/get_order_info",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "PayOrder",
			Router:           "/pay_order",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Token"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Token"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           "/login",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Token"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Token"],
		beego.ControllerComments{
			Method:           "Verify",
			Router:           "/verify",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"],
		beego.ControllerComments{
			Method:           "Earnlist",
			Router:           "/my/earnlist",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/my/index",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:User"],
		beego.ControllerComments{
			Method:           "OrderList",
			Router:           "/my/orders",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}

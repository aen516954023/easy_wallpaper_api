package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
			Router:           "/get_orders_all",
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
			Method:           "GetPhone",
			Router:           "/phone",
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

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Workers"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Workers"],
		beego.ControllerComments{
			Method:           "OrderTaking",
			Router:           "/order_taking",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}

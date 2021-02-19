package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
			Method:           "Uploads",
			Router:           "/uploads",
			AllowHTTPMethods: []string{"post"},
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
			Method:           "AdvanceOrder",
			Router:           "/advance_order",
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
			Method:           "ConfirmOrderChange",
			Router:           "/confirm_order_change",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "OrderDetails",
			Router:           "/get_order_details",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "Index",
			Router:           "/get_orders",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "OrderCancel",
			Router:           "/order_cancel",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "OrderPages",
			Router:           "/order_page",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"] = append(beego.GlobalControllerRouter["easy_wallpaper_api/controllers:Orders"],
		beego.ControllerComments{
			Method:           "PayAdvanceOrder",
			Router:           "/pay_advance_order",
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

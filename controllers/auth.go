package controllers

import (
	"easy_wallpaper_api/models"
	"easy_wallpaper_api/util"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
	"time"
)

// Operations about Users
type Auth struct {
	beego.Controller
}

var bm, _ = cache.NewCache("memory", `{"interval":60}`)

// @Title 邮箱注册 Create
// @Description 注册
// @Param	email		query 	string	true		"The email for auth"
// @Param	password		query 	string	true		"The password for auth"
// @Param	passveri		query 	string	true		"邮箱验证码"
// @Param	welcode		query 	string	true		"邀请码"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /EmailCreate [post]
func (this *Auth) EmailCreate() {
	//检测该用户是否授权完成
	// 获取数据
	email := this.GetString("email")
	password := this.GetString("password")
	passveri := this.GetString("passveri")
	welcode, _ := this.GetInt("welcode")
	//验证器
	// 校验数据
	// 判断是否为空
	if password == "" || email == "" || passveri == "" {
		this.Data["json"] = ReturnError(0, "数据填写不完整，请重新输入~")
		this.ServeJSON()
		this.StopRun()
	}
	//验证器
	errb := validation.Validate(welcode,
		validation.Required.Error("请填写正确邀请码！"),
	)
	if errb != nil {
		this.Data["json"] = ReturnError(0, "请填写正确邀请码!")
		this.ServeJSON()
		this.StopRun()
	}
	// 使用正则判断邮箱格式
	regex, _ := regexp.Compile("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}")
	resEmail := regex.FindString(email)
	if resEmail == "" {
		this.Data["json"] = ReturnError(0, "邮箱格式不正确，请重新输入~")
		this.ServeJSON()
		this.StopRun()
	}
	//判断邮箱验证码是否正确
	if bm.IsExist(email) {
		if bm.Get(email) != passveri {
			this.Data["json"] = ReturnError(0, "邮件验证码错误，请重新核对！")
			this.ServeJSON()
			this.StopRun()
		}
	} else {
		this.Data["json"] = ReturnError(0, "邮件验证码失效，请重新获取！")
		this.ServeJSON()
		this.StopRun()
	}
	status, err := models.CreateUser(email, password, welcode)
	if status {
		if err == nil {
			//生成token并返回
			user := util.User{Member: email}
			token, err := util.GenerateToken(&user, 0)
			if err == nil {
				msg := map[string]string{
					"token": token,
				}
				this.Data["json"] = ReturnSuccess(1, "success", msg, 1)
				this.ServeJSON()
			} else {
				this.Data["json"] = ReturnError(0, "生成Token失败")
				this.ServeJSON()
			}
		} else {
			//插入数据库失败
			this.Data["json"] = ReturnError(0, "参数错误！")
			this.ServeJSON()
		}

	} else {
		//邀请码无效
		this.Data["json"] = ReturnError(0, "注册失败，该用户已注册！")
	}
	this.ServeJSON()
	//}

}

// @Title 发送邮箱验证码 SendEmailVer
// @Description 发送邮箱验证码
// @Param	email		query 	string	true		"The email for auth"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /SendEmailVer [post]
func (this *Auth) SendEmailVer() {
	email := this.GetString("email")
	regex, _ := regexp.Compile("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}")
	resEmail := regex.FindString(email)
	if resEmail == "" {
		this.Data["json"] = ReturnError(40001, "邮箱格式不正确，请重新输入~")
		this.ServeJSON()
		this.StopRun()
	}
	// 发送激活邮件
	emailConfig := `{"username":"1017093063@qq.com","password":"cugnefgflpzsbcag","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "1017093063@qq.com" // 指定发件人的邮箱地址
	emailConn.To = []string{email}       // 指定收件人邮箱地址
	emailConn.Subject = "用户激活"           // 指定邮件的标题
	code := CreateCaptcha()
	// 发给用户的是激活请求地址
	//emailConn.HTML = `尊敬的` + user.Name + `，您好<br><br>感谢您注册，为了避免您忘记账号或密码导致您的账户无法找回，请您验证Email地址。<br><br>请复制粘贴下面的链接至浏览器地址栏打开：<br><br>127.0.0.1:8080/active?id=` + strconv.Itoa(user.Id) + `<br>`
	emailConn.HTML = `Welcome` + email + `<br>Your  register  verification code is:` + code + `,<br>Thanks for your support!You have received this email because you are using the computing power sharing platform`
	err := emailConn.Send()
	if err != nil {
		this.Data["json"] = ReturnError(40001, "发送激活邮件失败，请检查格式是否正确，重新发送~")
		this.ServeJSON()
		this.StopRun()
	} else {
		//将code储存到缓存中
		fmt.Printf("code", code)
		bm.Put(email, code, 120*time.Second)
		this.Data["json"] = ReturnSuccess(1, "success", "发送邮件成功", 1)
		this.ServeJSON()
		this.StopRun()
	}
}

// @Title 邮箱登录 EmailLogin
// @Description 邮箱登录
// @Param	email		query 	string	true		"The email for auth"
// @Param	password		query 	string	true		"The password for auth"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /EmailLogin [post]
func (this *Auth) EmailLogin() {
	// 获取数据
	Email := this.GetString("email")
	password := this.GetString("password")
	// 校验数据
	if Email == "" || password == "" {
		this.Data["json"] = ReturnError(0, "邮箱或密码不能为空")
		this.ServeJSON()
		this.StopRun()
	}
	regex, _ := regexp.Compile("\\w[-\\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\\.)+[A-Za-z]{2,14}")
	resEmail := regex.FindString(Email)
	if resEmail == "" {
		this.Data["json"] = ReturnError(0, "邮箱格式不正确，请重新输入~")
		this.ServeJSON()
		this.StopRun()
	}
	status := models.EmailLogin(Email, password)
	if status {
		//登录成功，生成token信息
		//生成token并返回
		user := util.User{Member: Email}
		token, err := util.GenerateToken(&user, 0)
		if err == nil {
			msg := map[string]string{
				"token": token,
			}
			this.Data["json"] = ReturnSuccess(1, "success", msg, 1)
		} else {
			this.Data["json"] = ReturnError(0, "生成Token失败")
		}
		this.ServeJSON()
		this.StopRun()
	} else {
		this.Data["json"] = ReturnError(0, "账号密码错误，请核对后重新输入")
		this.ServeJSON()
		this.StopRun()
	}

}
func (this *Auth) Payc() {
	postdata := make(map[string]interface{})
	postdata["orders_code"] = "123321123" //订单号
	postdata["order_total"] = "1"         //支付总金额
	postdata["currency_code"] = "USD"     //币种，例：美金USD
	postdata["order_total_usd"] = "1"     //总折算美金金额
	postdata["notify_url"] = "baidu.com"  //支付结果回调地址
	postdata["products_id"] = "123"       //产品id
	postdata["products_name"] = "hahaha"  //产品名称
	postdata["products_price"] = "1"      //产品价格
	postdata["products_price_usd"] = "1"  //产品折算美金价格
	data := GetOrderUrl(postdata)
	fmt.Println("resc", data)
}

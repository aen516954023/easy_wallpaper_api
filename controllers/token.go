package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"easy_wallpaper_api/models"
	"easy_wallpaper_api/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"time"
)

//var bmToken, _ = cache.NewCache("memory", `{"interval":60}`)

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
	// 将sessionKey存储
	bmErr := bm.Put("session_key", res.SessionKey, 7200*time.Second)
	fmt.Println("bm", bmErr)
	// 查看用户是否存在 如果存在 刷新token 如果不存在 新增用户信息
	info, infoErr := models.GetMemberInfo(res.OpenId)
	fmt.Println(infoErr, info.Id)
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

//{"phoneNumber":"15938755991","purePhoneNumber":"15938755991","countryCode":"86","watermark":{"timestamp":1612099347,"appid":"wx7cfe2b3493f5cbc6"}} <nil>
type Watermark struct {
	Timestamp int64  `json"timestamp"`
	Appid     string `json"appid"`
}
type Phone struct {
	PhoneNumber     string    `json"phoneNumber"`
	PurePhoneNumber string    `json"purePhoneNumber"`
	CountryCode     string    `json"countryCode"`
	Watermark       Watermark `json"watermark"`
}

// @Title 获取手机信息
// @Description 获取加密的微信手机号信息
// @Param	encryptedData	query 	string	true	"加密信息"
// @Param	iv	query 	string	true	"加密算法的初始向量"
// @Success 200 {string} auth success
// @Failure 403 user not exist
// @router /phone [post]
func (this *Token) GetPhone() {
	token := this.Ctx.Request.Header.Get("token")
	if token == "" {
		this.Data["json"] = ReturnError(40000, "获取信息失败,token不能为空")
		this.ServeJSON()
		return
	}
	tokenInfo, tokenErr := util.ValidateToken(token)
	if tokenErr != nil || tokenInfo.Id == 0 {
		this.Data["json"] = ReturnError(40001, "获取信息失败,无效的TOKEN")
		this.ServeJSON()
		return
	}
	////sessionKey := "dG7DpbFgaUj0JeMgFEtrTA=="
	sessionKey := bm.Get("session_key")
	encryptedData := this.GetString("encryptedData")
	iv := this.GetString("iv")
	////fmt.Println(encryptedData,iv,sessionKey)
	str, err := Dncrypt(encryptedData, sessionKey.(string), iv)
	//fmt.Println(src,err)
	//str := `{"phoneNumber":"15938755991","purePhoneNumber":"15938755991","countryCode":"86","watermark":{"timestamp":1612099347,"appid":"wx7cfe2b3493f5cbc6"}}`
	var data Phone
	if err == nil {
		json.Unmarshal([]byte(str), &data)
		if data.Watermark.Appid != beego.AppConfig.String("appId") {
			this.Data["json"] = ReturnError(41003, "获取信息失败")
			this.ServeJSON()
			return
		}
		// 更新用户手机号码到数据表中
		if models.UpdatePhone(tokenInfo.Id, data.PurePhoneNumber) {
			this.Data["json"] = ReturnSuccess(0, "success", "", 0)
			this.ServeJSON()
			return
		}
	}
	this.Data["json"] = ReturnError(40000, "获取信息失败")
	this.ServeJSON()
}

// CBC 模式
//解密
/**
* rawData 原始加密数据
* key  密钥
* iv  向量
 */
func Dncrypt(rawData, key, iv string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	key_b, err_1 := base64.StdEncoding.DecodeString(key)
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if err_1 != nil {
		return "", err_1
	}
	dnData, err := AesCBCDncrypt(data, key_b, iv_b)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

// 解密
func AesCBCDncrypt(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

//去除填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

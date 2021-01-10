package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	orm.RegisterModel(new(User))
}

type User struct {
	Id            int
	Name          string `orm:"size(20);unique"` // 用户名
	Password      string `orm:"size(20)"`        // 密码
	Active        bool   `orm:"default(false)"`  // 是否激活，默认未激活
	Power         int    `orm:"default(0)"`      // 权限设置，0表示普通用户，1表示管理员
	Welcode       int    `orm:"default(0)"`      //
	Ppath         string `orm:"default(0)"`      //
	Pid           int
	Btc           float64
	Address       []*Address   `orm:"reverse(many)"`
	OrderInfo     []*OrderInfo `orm:"reverse(many)"`
	LastLoginTime int64
}

func CreateUser(Name string, Password string, welcode int) (bool, error) {
	//查询邀请码是否存在
	o := orm.NewOrm()
	cuser := new(User)
	err := o.QueryTable("user").Filter("welcode", welcode).One(cuser)
	if err != nil {
		fmt.Printf("cw", err)
		return false, err
	}
	user := User{}
	user.Name = Name
	//查询该用户是否已经注册
	err = o.Read(&user, "Name")
	if err == nil {
		return false, err
	} else {
		md5String := fmt.Sprintf("%x", md5.Sum([]byte(Password)))
		user.Password = md5String
		user.Welcode = generateRandomNumber(100000, 999999)
		user.Pid = cuser.Id
		user.LastLoginTime = time.Now().Unix()
		user.Ppath = cuser.Ppath + strconv.Itoa(cuser.Id) + ","
		_, err = o.Insert(&user)
		if err != nil {
			return false, err
		} else {
			return true, err
		}
	}
}
func FindUserByUsername(member string) (*User, error) {
	user := new(User)
	db := orm.NewOrm()
	err := db.QueryTable("user").Filter("name", member).One(user)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		fmt.Printf("Returned Multi Rows Not One")
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		fmt.Printf("Not row found")
	}
	return user, err
}
func EmailLogin(Email string, Password string) bool {
	o := orm.NewOrm()
	var user User
	user.Name = Email
	err := o.Read(&user, "Name")
	if err != nil {
		return false
	} else {
		md5String := fmt.Sprintf("%x", md5.Sum([]byte(Password)))
		if user.Password != md5String {
			return false
		} else {
			//更新登录时间
			user.LastLoginTime = time.Now().Unix()
			o.Update(&user)
			//登录成功
			return true
		}
	}

}

//增加收益
func AddUserFee(uid int, fee float64, ysy float64) error {
	o := orm.NewOrm()
	user := User{Id: uid}
	if o.Read(&user) == nil {
		user.Btc = fee + ysy
		if _, err := o.Update(&user); err == nil {
			//添加收益记录
			str := strconv.FormatFloat(fee, 'f', -1, 64)
			des := "获得：" + str + "BTC"
			AddRecord(uid, fee, des)
			return err
		} else {
			return err
		}
	} else {
		return o.Read(&user)
	}
}

//生成唯一邀请码
func generateRandomNumber(start int, end int) int {
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//生成随机数
	num := r.Intn((end - start)) + start
	return num
}

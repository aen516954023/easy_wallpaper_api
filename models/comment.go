package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EComments))
}

type EComments struct {
	Id          int
	MId         int
	OrderId     int
	Content     string
	Rate        int
	Image       string
	IsAnonymous int
	CreateAt    string
}

//查询评论
func GetUserComment(mId, oId int) (EComments, error) {
	o := orm.NewOrm()
	var data EComments
	err := o.QueryTable("e_comments").Filter("m_id", mId).Filter("order_id", oId).One(&data)
	return data, err
}

//添加评论
func AddComment(mId, oId, rate, isAnonymous int, content, image string) (bool, error) {
	o := orm.NewOrm()
	var data EComments
	data.MId = mId
	data.OrderId = oId
	data.Content = content
	data.Rate = rate
	data.Image = image
	data.IsAnonymous = isAnonymous
	data.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	insertId, err := o.Insert(&data)
	if err == nil && insertId > 0 {
		return true, err
	}
	return false, err
}

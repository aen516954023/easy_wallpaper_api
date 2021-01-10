package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

func init() {
	orm.RegisterModel(new(EarningRecord))
}

type EarningRecord struct {
	Id      int
	Uid     int
	Btc     float64
	Des     string
	AddTime int64
}

//增加记录
func AddRecord(uid int, btc float64, des string) (error, bool) {
	now := time.Now()
	currentTime := now.Unix()
	CreateAt := currentTime
	o := orm.NewOrm()
	var record EarningRecord
	record.Btc = btc
	record.Uid = uid
	record.Des = des
	record.AddTime = CreateAt
	_, err := o.Insert(&record)
	if err == nil {
		return err, true
	} else {
		return err, false
	}
}

//获取过去7天的收益
func GetcountryCapital(uid int) ([]map[string]interface{}, bool, error) {
	countryCapitalMap := make([]map[string]interface{}, 7)
	//查询最新的7条记录
	o := orm.NewOrm()
	var EarningRecord []EarningRecord
	_, err := o.Raw("SELECT btc,add_time FROM earning_record WHERE uid = ? order by add_time desc limit 7", uid).QueryRows(&EarningRecord)
	//for i := 0; i < 7; i++ {
	//	if countryCapitalMap[i] == nil {
	//		countryCapitalMap[i]["TM"] = 0
	//		countryCapitalMap[i]["PV"] = 0
	//	}
	//}
	//for key, value := range EarningRecord {
	//	countryCapitalMap[key]["TM"] = ParseTimeToString(value.AddTime)
	//	countryCapitalMap[key]["PV"] = value.Btc
	//}
	if len(EarningRecord) > 0 {
		for i := 0; i < len(EarningRecord); i++ {
			if countryCapitalMap[i] == nil {
				countryCapitalMap[i] = map[string]interface{}{}
			}
			countryCapitalMap[i]["TM"] = ParseTimeToString(EarningRecord[i].AddTime)
			countryCapitalMap[i]["PV"] = EarningRecord[i].Btc
		}
		return countryCapitalMap, true, err
	}
	return countryCapitalMap, false, err

	//for key, value := range EarningRecord {
	//	countryCapitalMap[key]["TM"] = ParseTimeToString(value.AddTime)
	//	countryCapitalMap[key]["PV"] = value.Btc
	//}

}
func GetEarnList(uid int) (int64, []EarningRecord, error) {
	o := orm.NewOrm()
	var data []EarningRecord
	num, err := o.QueryTable(new(EarningRecord)).Filter("uid", uid).All(&data)
	return num, data, err
}

//将int时间戳变成string
func ParseTimeToString(s int64) string {
	//fmt.Println(s)
	//fmt.Println(int64(s))
	tm := time.Unix(s, 0)
	return tm.Format("01/02")
}

package api

import (
	"fmt"
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of ap
 */

func HandleCreateWlan(new TblWlan) error {
	fmt.Println(new.Ssid)
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateWlan(id uint64, value map[string]interface{}) error {
	where := TblWlan{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteWlan(id uint64) error {
	where := TblWlan{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindWlanById(id uint64) (result TblWlan, err error) {
	where := TblWlan{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblWlan)
	return result, err
}
func HandleFindWlanByGroupUrl(groupUrl string) (result []TblWlan, err error) {
	where := TblWlan{GroupUrl:groupUrl}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblWlan)
	return result, err
}
func HandleFindWlanByGroupIdAndSsid(groupId uint64, ssid string) (result []TblWlan, err error) {
	where := TblWlan{GroupId:groupId,Ssid:ssid}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblWlan)
	return result, err
}
func HandleFindWlansByUrl(url string) (result []TblWlan, err error) {
	r, err := mysql_driver.Handle.FindLike(TblWlan{}, "url", url)
	result = r.([]TblWlan)
	return result, err
}
func HandleFindWlansAll() (result []TblWlan, err error) {
	r, err := mysql_driver.Handle.FindAll(TblWlan{})
	result = r.([]TblWlan)
	return result, err
}
func HandleFindWlansCountByGroupId(groupId uint64) (count int, err error) {
	where := TblWlan{GroupId:groupId}
	count, err = mysql_driver.Handle.FindAllCount(where)
	return count, err
}
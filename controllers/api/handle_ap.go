package api

import (
	"fmt"
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of ap
 */

func HandleCreateAp(new TblDevInfo) error {
	fmt.Println(new.Mac)
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateAp(id uint64, value map[string]interface{}) error {
	where := TblDevInfo{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleUpdateApByGroupUrl(groupUrl string, value map[string]interface{}) error {
	where := TblDevInfo{GroupUrl:groupUrl}
	return mysql_driver.Handle.Update(where, value)
}
func HandleUpdateApByMac(mac string, value map[string]interface{}) error {
	where := TblDevInfo{Mac:mac}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteAp(id uint64) error {
	where := TblDevInfo{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindApById(id uint64) (result TblDevInfo, err error) {
	where := TblDevInfo{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblDevInfo)
	return result, err
}
func HandleFindApByMac(mac string) (result TblDevInfo, err error) {
	where := TblDevInfo{Mac:mac}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblDevInfo)
	return result, err
}
func HandleFindApsByUrl(url string) (aps []string, err error) {
	r, err := mysql_driver.Handle.FindAll(TblDevInfo{GroupUrl:url})
	result := r.([]TblDevInfo)
	for i := 0; i < len(result); i++ {
		aps = append(aps, result[i].Mac)
	}
	return aps, err
}
func HandleFindApsAll() (result []TblDevInfo, err error) {
	r, err := mysql_driver.Handle.FindAll(TblDevInfo{})
	result = r.([]TblDevInfo)
	return result, err
}
func HandleFindApCountByGroupId(groupId uint64) (all int, online int, err error) {
	where := TblDevInfo{GroupId:groupId}
	all, err = mysql_driver.Handle.FindAllCount(where)
	if all <= 0 {
		online = 0
	} else {
		wh := TblDevStat{GroupId:groupId, State:"online"}
		online, err = mysql_driver.Handle.FindAllCount(wh)
	}
	return all, online, err
}
func HandleFindApStateByMac(mac string) (result TblDevStat, err error) {
	where := TblDevStat{Mac:mac}
	r, err := mysql_driver.Handle.FindFirst(where)
	if err != nil {
		return result, err
	}
	result = r.(TblDevStat)
	return result, nil
}
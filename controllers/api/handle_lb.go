package api

import (
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of load balance
 */

func HandleCreateLb(new TblLoadBalanceGroup) error {
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateLb(id uint64, value map[string]interface{}) error {
	where := TblLoadBalanceGroup{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteLb(id uint64) error {
	where := TblLoadBalanceGroup{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindLbById(id uint64) (result TblLoadBalanceGroup, err error) {
	where := TblLoadBalanceGroup{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblLoadBalanceGroup)
	return result, err
}
func HandleFindLbByGroupUrl(groupUrl string) (result []TblLoadBalanceGroup, err error) {
	where := TblLoadBalanceGroup{GroupUrl:groupUrl}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblLoadBalanceGroup)
	return result, err
}
func HandleFindLbsAll() (result []TblLoadBalanceGroup, err error) {
	r, err := mysql_driver.Handle.FindAll(TblLoadBalanceGroup{})
	result = r.([]TblLoadBalanceGroup)
	return result, err
}
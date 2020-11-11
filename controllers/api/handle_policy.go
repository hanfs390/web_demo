package api

import (
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of policy
 */

func HandleCreatePolicy(new TblPolicy) error {
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdatePolicy(id uint64, value map[string]interface{}) error {
	where := TblPolicy{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeletePolicy(id uint64) error {
	where := TblPolicy{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindPolicyById(id uint64) (result TblPolicy, err error) {
	where := TblPolicy{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblPolicy)
	return result, err
}
func HandleFindPolicyByGroupUrl(groupUrl string) (result []TblPolicy, err error) {
	where := TblPolicy{GroupUrl:groupUrl}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblPolicy)
	return result, err
}
func HandleFindPolicysByUrl(url string) (result []TblPolicy, err error) {
	r, err := mysql_driver.Handle.FindLike(TblPolicy{}, "url", url)
	result = r.([]TblPolicy)
	return result, err
}
func HandleFindPolicysAll() (result []TblPolicy, err error) {
	r, err := mysql_driver.Handle.FindAll(TblPolicy{})
	result = r.([]TblPolicy)
	return result, err
}
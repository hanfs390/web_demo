package api

import (
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of black/white policy
 */

func HandleCreateBWlist(new TblBWList) error {
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateBWlist(id uint64, value map[string]interface{}) error {
	where := TblBWList{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteBWlist(id uint64) error {
	where := TblBWList{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindBWlistById(id uint64) (result TblBWList, err error) {
	where := TblBWList{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblBWList)
	return result, err
}
func HandleFindBWlistByGroupUrl(groupUrl string) (result []TblBWList, err error) {
	where := TblBWList{GroupUrl:groupUrl}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblBWList)
	return result, err
}
func HandleFindBWlistsByUrl(url string) (result []TblBWList, err error) {
	r, err := mysql_driver.Handle.FindLike(TblBWList{}, "url", url)
	result = r.([]TblBWList)
	return result, err
}
func HandleFindBWlistsAll() (result []TblBWList, err error) {
	r, err := mysql_driver.Handle.FindAll(TblBWList{})
	result = r.([]TblBWList)
	return result, err
}
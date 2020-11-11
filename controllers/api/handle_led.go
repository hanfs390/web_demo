package api

import (
	"ByzoroAC/models/mysql_driver"
	. "ByzoroAC/common"
)

/**
 * the cmd of led
 */

func HandleCreateLed(new TblLED) error {
	return mysql_driver.Handle.Insert(new)
}
func HandleUpdateLed(id uint64, value map[string]interface{}) error {
	where := TblLED{Id:id}
	return mysql_driver.Handle.Update(where, value)
}
func HandleDeleteLed(id uint64) error {
	where := TblLED{Id:id}
	return mysql_driver.Handle.Delete(where)
}
func HandleFindLedById(id uint64) (result TblLED, err error) {
	where := TblLED{Id:id}
	r, err := mysql_driver.Handle.FindFirst(where)
	result = r.(TblLED)
	return result, err
}
func HandleFindLedByGroupUrl(groupUrl string) (result []TblLED, err error) {
	where := TblLED{GroupUrl:groupUrl}
	r, err := mysql_driver.Handle.FindAll(where)
	result = r.([]TblLED)
	return result, err
}
func HandleFindLedsByUrl(url string) (result []TblLED, err error) {
	r, err := mysql_driver.Handle.FindLike(TblLED{}, "url", url)
	result = r.([]TblLED)
	return result, err
}
func HandleFindLedsAll() (result []TblLED, err error) {
	r, err := mysql_driver.Handle.FindAll(TblLED{})
	result = r.([]TblLED)
	return result, err
}